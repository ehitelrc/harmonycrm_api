-- +goose Up
-- +goose StatementBegin
--
-- PostgreSQL database dump
--

\restrict rbGpHnXKUrGkQ3B0vN2jGJhpUDbrfqkscoYa29rZTivHukUvT1H24zSh6OIfMh1

-- Dumped from database version 16.13 (Ubuntu 16.13-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.13

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: -
--

-- *not* creating schema, since initdb creates it


--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: fn_update_case_payment_found(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.fn_update_case_payment_found() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Solo actualiza si viene un case_id válido
    IF NEW.case_id IS NOT NULL THEN
        UPDATE cases
        SET payment_found = true,
            updated_at = now()
        WHERE id = NEW.case_id;
    END IF;

    RETURN NEW;
END;
$$;


--
-- Name: fn_update_erp_payment_confirmation(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.fn_update_erp_payment_confirmation() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN

    UPDATE public.erp_payment_confirmation
    SET harmony_state = 'recibido_por_harmony',
        updated_at = now()
    WHERE reference_number = NEW.reference_number;

    RETURN NEW;

END;
$$;


--
-- Name: fn_validate_user_company(integer, text, text); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.fn_validate_user_company(p_company_id integer, p_email text, p_password text) RETURNS TABLE(has_permission boolean, user_id integer)
    LANGUAGE plpgsql
    AS $_$
DECLARE
    v_user_id       INTEGER;
    v_password_hash TEXT;
    v_has_access    BOOLEAN;
BEGIN
    -- Buscar el usuario y su hash
    SELECT u.id, u.password_hash
      INTO v_user_id, v_password_hash
      FROM public.users u
     WHERE u.email = p_email
     LIMIT 1;

    IF v_user_id IS NULL THEN
        -- no existe usuario
        has_permission := FALSE;
        user_id        := NULL;
        RETURN NEXT;
        RETURN;
    END IF;

    -- (Opcional) Validación de formato de hash para evitar invalid salt
    IF v_password_hash IS NULL
       OR v_password_hash = ''
       OR v_password_hash !~ '^\$2[aby]\$[0-9]{2}\$[./0-9A-Za-z]{53}$' THEN
        has_permission := FALSE;
        user_id        := NULL;
        RETURN NEXT;
        RETURN;
    END IF;

    -- Validar password
    IF crypt(p_password, v_password_hash) <> v_password_hash THEN
        has_permission := FALSE;
        user_id        := NULL;
        RETURN NEXT;
        RETURN;
    END IF;

    -- Verificar acceso a la compañía (calificando columnas con alias)
    SELECT EXISTS (
        SELECT 1
          FROM public.user_company_roles AS ucr
         WHERE ucr.user_id   = v_user_id
           AND ucr.company_id = p_company_id
    ) INTO v_has_access;

    has_permission := v_has_access;
    user_id        := CASE WHEN v_has_access THEN v_user_id ELSE NULL END;

    RETURN NEXT;
END;
$_$;


--
-- Name: get_complete_trip_details(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_complete_trip_details() RETURNS json
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN (
        SELECT json_agg(
            json_build_object(
                -- TRIP INFO
                'trip_id', t.trip_id,
                'status', t.status,
                'rate_type', t.rate_type,
                'comments', t.comments,
                'created_at', t.created_at,
                'updated_at', t.updated_at,
                'scheduled_date', t.scheduled_date,
                'assigned_date', t.assigned_date,
                'started_date', t.started_date,
                'finished_date', t.finished_date,

                -- Duraciones en minutos
                'minutes_since_created', 
                    ROUND(EXTRACT(EPOCH FROM now() - t.created_at) / 60),
                'minutes_since_started', 
                    CASE 
                        WHEN t.started_date IS NOT NULL 
                        THEN ROUND(EXTRACT(EPOCH FROM now() - t.started_date) / 60)
                        ELSE NULL 
                    END,

                -- COMPANY INFO
                'company_id', c.company_id,
                'company_name', c.company_name,
                'company_phone', c.company_phone,
                'company_mail', c.company_mail,
                'billing_model', c.billing_model,
                'community_company', c.community_company,

                -- BRANCH INFO
                'branch_id', b.branch_id,
                'branch_name', b.branch_name,
                'branch_phone', b.branch_phone,
                'branch_mail', b.branch_mail_,
                'branch_address', b.address,

                -- VEHICLE INFO
                'vehicle_id', v.vehicle_id,
                'vehicle_model', v.model,
                'vehicle_year', v.vehicle_year,
                'plate_number', v.plate_number,
                'vehicle_color', v.color,
                'vehicle_style', v.style,
                'engine_number', v.engine_number,
                'chassis_number', v.chassis_number,
                'make_name', vm.make_name,
                'profile_name', vp.profile_name,
                'profile_description', vp.description,

                -- DRIVER INFO
                'driver_id', d.user_id,
                'driver_name', CONCAT(d.first_name, ' ', d.last_name),
                'driver_phone', d.phone_number,
                'driver_picture', d.profile_picture,

                -- OWNER INFO (USER)
                'user_id', u.user_id,
                'user_name', CONCAT(u.first_name, ' ', u.last_name),
                'user_phone', u.phone_number,
                'user_email', u.email,

                -- INSURANCE INFO
                'insurance_policy', CASE WHEN t.requires_insurance = true AND t.insurance_id IS NOT NULL THEN (
                    SELECT json_build_object(
                        'policy_number', il.policy_number,
                        'provider_name', ip.provider_name
                    )
                    FROM insurance_rates ir
                    LEFT JOIN insurance_list il ON ir.insurance_list_id = il.insurance_list_id
                    LEFT JOIN insurance_providers ip ON il.provider_id = ip.provider_id
                    WHERE ir.rate_id = t.insurance_id
                ) ELSE NULL END,

                -- LAST VEHICLE LOCATION
                'last_vehicle_location', (
                    SELECT json_build_object(
                        'latitude', tvl.latitude,
                        'longitude', tvl.longitude,
                        'capture_timestamp', tvl.capture_timestamp
                    )
                    FROM trip_vehicle_locations tvl
                    WHERE tvl.trip_id = t.trip_id
                    ORDER BY tvl.capture_timestamp DESC
                    LIMIT 1
                ),
                
                'has_photos', EXISTS(
                    SELECT 1 FROM trip_images ti WHERE ti.trip_id = t.trip_id
                ),

                -- TRIP LOCATIONS
                'trip_locations', (
                    SELECT json_agg(
                        json_build_object(
                            'location_id', tl.location_id,
                            'sequence', tl.sequence,
                            'is_origin', tl.is_origin,
                            'name', tl.name,
                            'full_address', tl.full_address,
                            'latitude', tl.latitude,
                            'longitude', tl.longitude,
                            'status', tl.status,
                            'start_time', tl.start_time,
                            'end_time', tl.end_time
                        )
                        ORDER BY tl.sequence
                    )
                    FROM trip_locations tl
                    WHERE tl.trip_id = t.trip_id
                )
            )
            ORDER BY t.trip_id
        )
        FROM trips t
        LEFT JOIN companies c ON c.company_id = t.company_id
        LEFT JOIN branches b ON b.branch_id = t.branch_id
        LEFT JOIN vehicles v ON v.vehicle_id = t.vehicle_id
        LEFT JOIN vehicle_makes vm ON vm.make_id = v.make_id
        LEFT JOIN vehicle_profiles vp ON vp.profile_id = v.profile_id
        LEFT JOIN users d ON d.user_id = t.driver_id
        LEFT JOIN users u ON u.user_id = t.user_id
        WHERE t.is_closed = false
    );
END;
$$;


--
-- Name: get_complete_trip_details_by_company(bigint); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_complete_trip_details_by_company(company_id_param bigint) RETURNS json
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN (
        SELECT json_agg(
            json_build_object(
                -- TRIP INFO
                'trip_id', t.trip_id,
                'status', t.status,
                'rate_type', t.rate_type,
                'comments', t.comments,
                'created_at', t.created_at,
                'updated_at', t.updated_at,
                'scheduled_date', t.scheduled_date,
                'assigned_date', t.assigned_date,
                'started_date', t.started_date,
                'finished_date', t.finished_date,

                -- Duraciones en minutos
                'minutes_since_created', 
                    ROUND(EXTRACT(EPOCH FROM now() - t.created_at) / 60),
                'minutes_since_started', 
                    CASE 
                        WHEN t.started_date IS NOT NULL 
                        THEN ROUND(EXTRACT(EPOCH FROM now() - t.started_date) / 60)
                        ELSE NULL 
                    END,

                -- COMPANY INFO
                'company_id', c.company_id,
                'company_name', c.company_name,
                'company_phone', c.company_phone,
                'company_mail', c.company_mail,
                'billing_model', c.billing_model,
                'community_company', c.community_company,

                -- BRANCH INFO
                'branch_id', b.branch_id,
                'branch_name', b.branch_name,
                'branch_phone', b.branch_phone,
                'branch_mail', b.branch_mail_,
                'branch_address', b.address,

                -- VEHICLE INFO
                'vehicle_id', v.vehicle_id,
                'vehicle_model', v.model,
                'vehicle_year', v.vehicle_year,
                'plate_number', v.plate_number,
                'vehicle_color', v.color,
                'vehicle_style', v.style,
                'engine_number', v.engine_number,
                'chassis_number', v.chassis_number,
                'make_name', vm.make_name,
                'profile_name', vp.profile_name,
                'profile_description', vp.description,

                -- DRIVER INFO
                'driver_id', d.user_id,
                'driver_name', CONCAT(d.first_name, ' ', d.last_name),
                'driver_phone', d.phone_number,
                'driver_picture', d.profile_picture,

                -- OWNER INFO (USER)
                'user_id', u.user_id,
                'user_name', CONCAT(u.first_name, ' ', u.last_name),
                'user_phone', u.phone_number,
                'user_email', u.email,

                -- INSURANCE INFO
                'insurance_policy', CASE WHEN t.requires_insurance = true AND t.insurance_id IS NOT NULL THEN (
                    SELECT json_build_object(
                        'policy_number', il.policy_number,
                        'provider_name', ip.provider_name
                    )
                    FROM insurance_rates ir
                    LEFT JOIN insurance_list il ON ir.insurance_list_id = il.insurance_list_id
                    LEFT JOIN insurance_providers ip ON il.provider_id = ip.provider_id
                    WHERE ir.rate_id = t.insurance_id
                ) ELSE NULL END,

                -- LAST VEHICLE LOCATION
                'last_vehicle_location', (
                    SELECT json_build_object(
                        'latitude', tvl.latitude,
                        'longitude', tvl.longitude,
                        'capture_timestamp', tvl.capture_timestamp
                    )
                    FROM trip_vehicle_locations tvl
                    WHERE tvl.trip_id = t.trip_id
                    ORDER BY tvl.capture_timestamp DESC
                    LIMIT 1
                ),

                'has_photos', EXISTS(
                    SELECT 1 FROM trip_images ti WHERE ti.trip_id = t.trip_id
                ),

                -- TRIP LOCATIONS
                'trip_locations', (
                    SELECT json_agg(
                        json_build_object(
                            'location_id', tl.location_id,
                            'sequence', tl.sequence,
                            'is_origin', tl.is_origin,
                            'name', tl.name,
                            'full_address', tl.full_address,
                            'latitude', tl.latitude,
                            'longitude', tl.longitude,
                            'status', tl.status,
                            'start_time', tl.start_time,
                            'end_time', tl.end_time
                        )
                        ORDER BY tl.sequence
                    )
                    FROM trip_locations tl
                    WHERE tl.trip_id = t.trip_id
                )
            )
            ORDER BY t.trip_id
        )
        FROM trips t
        LEFT JOIN companies c ON c.company_id = t.company_id
        LEFT JOIN branches b ON b.branch_id = t.branch_id
        LEFT JOIN vehicles v ON v.vehicle_id = t.vehicle_id
        LEFT JOIN vehicle_makes vm ON vm.make_id = v.make_id
        LEFT JOIN vehicle_profiles vp ON vp.profile_id = v.profile_id
        LEFT JOIN users d ON d.user_id = t.driver_id
        LEFT JOIN users u ON u.user_id = t.user_id
        WHERE t.is_closed = false AND t.company_id = company_id_param
    );
END;
$$;


--
-- Name: get_driver_trips(integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_driver_trips(driver_id_filter integer) RETURNS json
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN (
        SELECT json_agg(trip_data ORDER BY trip_data->'trip'->>'trip_id')
        FROM (
            SELECT DISTINCT ON (t.trip_id)
                json_build_object(
                    'trip', json_build_object(
                        'trip_id', t.trip_id,
                        'company_id', t.company_id,
                        'branch_id', t.branch_id,
                        'customer_name', vut.customer_name,
                        'user_id', t.user_id,
                        'vehicle_id', t.vehicle_id,
                        'vehicle', json_build_object(
                            'make_id', v.make_id,
                            'make_name', m.make_name,
                            'model', v.model,
                            'vehicle_year', v.vehicle_year,
                            'plate_number', v.plate_number,
                            'engine_number', v.engine_number,
                            'chassis_number', v.chassis_number,
                            'color', v.color,
                            'style', v.style
                        ),
                        'driver_id', t.driver_id,
                        'driver_name', CONCAT(us.first_name, ' ', us.last_name),
                        'driver_phone', us.phone_number,
                        'driver_image', us.profile_picture,
                        'is_scheduled', t.is_scheduled,
                        'scheduled_date', t.scheduled_date,
                        'requires_insurance', t.requires_insurance,
                        'is_batch_service', t.is_batch_service,
                        'batch_identifier', t.batch_identifier,
                        'status', t.status,
                        'created_at', t.created_at,
                        'updated_at', t.updated_at,
                        'rate_type', t.rate_type,
                        'comments', t.comments,
                        'insurance_policy', CASE WHEN t.requires_insurance = true AND t.insurance_id IS NOT NULL THEN (
                            SELECT json_build_object(
                                'policy_number', il.policy_number,
                                'provider_name', ip.provider_name
                            )
                            FROM insurance_rates ir
                            LEFT JOIN insurance_list il ON ir.insurance_list_id = il.insurance_list_id
                            LEFT JOIN insurance_providers ip ON il.provider_id = ip.provider_id
                            WHERE ir.rate_id = t.insurance_id
                        ) ELSE NULL END,
                        'trip_vehicle_condition', (
                            SELECT json_build_object(
                                'driver_comments', tvc.driver_comments,
                                'has_engine_alert', tvc.has_engine_alert,
                                'has_engine_noises', tvc.has_engine_noises,
                                'fuel_level', tvc.fuel_level,
                                'has_internal_damages', tvc.has_internal_damages,
                                'odometer_reading', tvc.odometer_reading,
                                'has_damaged_night_lights', tvc.has_damaged_night_lights,
                                'has_scratches', tvc.has_scratches,
                                'has_damaged_security_lights', tvc.has_damaged_security_lights,
                                'has_damaged_tires', tvc.has_damaged_tires,
                                'has_damaged_windows', tvc.has_damaged_windows,
                                'is_clean_interior', tvc.is_clean_interior,
                                'is_clean_exterior', tvc.is_clean_exterior
                            )
                            FROM trip_vehicle_condition tvc
                            WHERE tvc.trip_id = t.trip_id
                        ),
                        'has_photos', EXISTS(
                            SELECT 1 FROM trip_images ti WHERE ti.trip_id = t.trip_id
                        ),
                        'trip_locations', (
                            SELECT json_agg(
                                json_build_object(
                                    'location_id', tl.location_id,
                                    'sequence', tl.sequence,
                                    'is_origin', tl.is_origin,
                                    'name', tl.name,
                                    'full_address', tl.full_address,
                                    'contact_phone', tl.contact_phone,
                                    'latitude', tl.latitude,
                                    'longitude', tl.longitude,
                                    'is_authorized_person', tl.is_authorized_person,
                                    'authorized_person_name', tl.authorized_person_name,
                                    'status', tl.status,
                                    'start_time', tl.start_time,
                                    'end_time', tl.end_time
                                )
                                ORDER BY tl.sequence
                            )
                            FROM trip_locations tl
                            WHERE tl.trip_id = t.trip_id
                        )
                    )
                ) AS trip_data
            FROM trips t
            LEFT JOIN vehicles v ON t.vehicle_id = v.vehicle_id
            LEFT JOIN public.vehicle_makes m ON v.make_id = m.make_id
            LEFT JOIN users us ON us.user_id = t.driver_id

          
            LEFT JOIN public.vw_user_company_type vut 
                ON vut.company_id = t.company_id
                AND vut.branch_id = t.branch_id
                AND vut.user_id = t.user_id
            WHERE 
				t.is_closed = false
				AND 
					(
						(t.status = 'ASSIGNED' AND t.driver_id = driver_id_filter)
						OR
						(t.status = 'STARTED' AND t.driver_id = driver_id_filter)
						OR 
						(
							(
								t.status = 'PENDING' AND t.driver_id IS NULL
								AND EXISTS (
								                SELECT 1
								                FROM preferred_drivers pd
								                WHERE pd.company_id = t.company_id
								                  AND pd.driver_id = driver_id_filter
								            )
							)
						)

						-- Podría fatlar otro estado y asignado al chofer
					)
			 
			    
				
             
            ORDER BY t.trip_id
        ) trips_unique
    );
END;
$$;


--
-- Name: get_last_driver_trip(integer, text); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_last_driver_trip(driver_id_filter integer, plate_number_filter text) RETURNS json
    LANGUAGE plpgsql STABLE COST 200
    AS $$
BEGIN
    RETURN (
        SELECT json_build_object(
            'trip_id', t.trip_id,
            'company_id', t.company_id,
            'branch_id', t.branch_id,
            'user_id', t.user_id,
			'customer_name', t.customer_name,
            'vehicle_id', t.vehicle_id,
            'vehicle', json_build_object(
                'make_id', t.make_id,
                'make_name', t.make_name,
                'model', t.model,
                'vehicle_year', t.vehicle_year,
                'plate_number', t.plate_number,
                'engine_number', t.engine_number,
                'chassis_number', t.chassis_number,
                'color', t.color,
                'style', t.style
            ),
            'driver_id', t.driver_id,
            'driver_name', CONCAT(t.first_name, ' ', t.last_name),
            'driver_phone', t.phone_number,
            'driver_image', t.profile_picture,
            'is_scheduled', t.is_scheduled,
            'scheduled_date', t.scheduled_date,
            'requires_insurance', t.requires_insurance,
            'is_batch_service', t.is_batch_service,
            'batch_identifier', t.batch_identifier,
            'status', t.status,
            'created_at', t.created_at,
            'updated_at', t.updated_at,
            'rate_type', t.rate_type,
            'comments', t.comments,
            'insurance_policy', CASE WHEN t.requires_insurance = true AND t.insurance_id IS NOT NULL THEN (
                SELECT json_build_object(
                    'policy_number', il.policy_number,
                    'provider_name', ip.provider_name
                )
                FROM insurance_rates ir
                LEFT JOIN insurance_list il ON ir.insurance_list_id = il.insurance_list_id
                LEFT JOIN insurance_providers ip ON il.provider_id = ip.provider_id
                WHERE ir.rate_id = t.insurance_id
            ) ELSE NULL END,
            'has_photos', EXISTS(
                SELECT 1 FROM trip_images ti WHERE ti.trip_id = t.trip_id
            )
        )
        FROM (
            SELECT t.*, vut.customer_name, v.make_id, m.make_name, v.model, v.vehicle_year, v.plate_number,
                   v.engine_number, v.chassis_number, v.color, v.style, us.last_name, us.first_name, us.phone_number, us.profile_picture
            FROM trips t
            LEFT JOIN vehicles v ON t.vehicle_id = v.vehicle_id
            LEFT JOIN users us ON us.user_id = t.driver_id
            LEFT JOIN public.vehicle_makes m ON v.make_id = m.make_id
			LEFT JOIN public.vw_user_company_type vut ON vut.company_id = t.company_id
			AND vut.branch_id = t.branch_id
			AND vut.user_id = t.user_id
		WHERE (
                t.driver_id = driver_id_filter
            )
            AND v.plate_number = plate_number_filter
            ORDER BY t.trip_id DESC
            LIMIT 1
        ) AS t
    );
END;
$$;


--
-- Name: get_trip_details(integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_trip_details(user_id_filter integer) RETURNS json
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN (
        SELECT json_agg(
            json_build_object(
                'trip', json_build_object(
                    'trip_id', t.trip_id,
                    'company_id', t.company_id,
                    'branch_id', t.branch_id,
					'customer_name', vut.customer_name,
                    'user_id', t.user_id,
                    'vehicle_id', t.vehicle_id,
                    'vehicle', json_build_object(
                        'make_id', v.make_id,
                        'make_name', m.make_name,
                        'model', v.model,
                        'vehicle_year', v.vehicle_year,
                        'plate_number', v.plate_number,
						'engine_number', v.engine_number,
						'chassis_number', v.chassis_number,
                        'color', v.color,
                        'style', v.style
                    ),
                    'driver_id', t.driver_id,
                    'driver_name', CONCAT(us.first_name, ' ', us.last_name),
                    'driver_phone', us.phone_number,
                    'driver_image', us.profile_picture,
                    'is_scheduled', t.is_scheduled,
                    'scheduled_date', t.scheduled_date,
                    'requires_insurance', t.requires_insurance,
                    'is_batch_service', t.is_batch_service,
                    'batch_identifier', t.batch_identifier,
                    'status', t.status,
                    'created_at', t.created_at,
                    'updated_at', t.updated_at,
                    'rate_type', t.rate_type,
                    'comments', t.comments,
                    'insurance_policy', CASE WHEN t.requires_insurance = true AND t.insurance_id IS NOT NULL THEN (
                        SELECT json_build_object(
                            'policy_number', il.policy_number,
                            'provider_name', ip.provider_name
                        )
                        FROM insurance_rates ir
                        LEFT JOIN insurance_list il ON ir.insurance_list_id = il.insurance_list_id
                        LEFT JOIN insurance_providers ip ON il.provider_id = ip.provider_id
                        WHERE ir.rate_id = t.insurance_id
                    ) ELSE NULL END,
                    'trip_vehicle_condition', (
                        SELECT json_build_object(
                            'driver_comments', tvc.driver_comments,
                            'has_engine_alert', tvc.has_engine_alert,
                            'has_engine_noises', tvc.has_engine_noises,
                            'fuel_level', tvc.fuel_level,
                            'has_internal_damages', tvc.has_internal_damages,
                            'odometer_reading', tvc.odometer_reading,
                            'has_damaged_night_lights', tvc.has_damaged_night_lights,
                            'has_scratches', tvc.has_scratches,
                            'has_damaged_security_lights', tvc.has_damaged_security_lights,
                            'has_damaged_tires', tvc.has_damaged_tires,
                            'has_damaged_windows', tvc.has_damaged_windows,
                            'is_clean_interior', tvc.is_clean_interior,
                            'is_clean_exterior', tvc.is_clean_exterior
                        )
                        FROM trip_vehicle_condition tvc
                        WHERE tvc.trip_id = t.trip_id
                    ),
                    'has_photos', EXISTS(
                        SELECT 1 FROM trip_images ti WHERE ti.trip_id = t.trip_id
                    ),
                    'trip_locations', (
                        SELECT json_agg(
                            json_build_object(
                                'location_id', tl.location_id,
                                'sequence', tl.sequence,
                                'is_origin', tl.is_origin,
                                'name', tl.name,
                                'full_address', tl.full_address,
                                'contact_phone', tl.contact_phone,
                                'latitude', tl.latitude,
                                'longitude', tl.longitude,
                                'is_authorized_person', tl.is_authorized_person,
                                'authorized_person_name', tl.authorized_person_name,
                                'status', tl.status,
                                'start_time', tl.start_time,
                                'end_time', tl.end_time
                            )
                            ORDER BY tl.sequence
                        )
                        FROM trip_locations tl
                        WHERE tl.trip_id = t.trip_id
                    )
                )
            ) 
            ORDER BY t.trip_id -- Ordenamos aquí los resultados por trip_id
        )
        FROM trips t
        LEFT JOIN vehicles v ON t.vehicle_id = v.vehicle_id
        LEFT JOIN users us ON us.user_id = t.driver_id
        LEFT JOIN public.vehicle_makes m ON v.make_id = m.make_id
        LEFT JOIN trip_vehicle_condition tvc ON t.trip_id = tvc.trip_id
		LEFT JOIN public.vw_user_company_type vut ON vut.company_id = t.company_id
		AND vut.branch_id = t.branch_id
		AND vut.user_id = t.user_id
        WHERE t.is_closed = false
        AND (user_id_filter IS NULL OR t.user_id = user_id_filter)
    );
END;
$$;


--
-- Name: set_updated_at_column(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
    BEGIN
      NEW.updated_at = NOW();
      RETURN NEW;
    END;
    $$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: agent_department_assignments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.agent_department_assignments (
    id integer NOT NULL,
    agent_id integer,
    department_id integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: agent_department_assignments_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.agent_department_assignments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: agent_department_assignments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.agent_department_assignments_id_seq OWNED BY public.agent_department_assignments.id;


--
-- Name: agents; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.agents (
    user_id integer NOT NULL
);


--
-- Name: campaign_whatsapp_push; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.campaign_whatsapp_push (
    id bigint NOT NULL,
    campaign_id bigint NOT NULL,
    description character varying(50) NOT NULL,
    template_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    funnel_stage_id bigint,
    changed_by bigint,
    department_id bigint,
    channel_integration_id bigint
);


--
-- Name: campaign_whatsapp_push_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.campaign_whatsapp_push_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: campaign_whatsapp_push_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.campaign_whatsapp_push_id_seq OWNED BY public.campaign_whatsapp_push.id;


--
-- Name: campaign_whatsapp_push_leads; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.campaign_whatsapp_push_leads (
    id bigint NOT NULL,
    push_id bigint NOT NULL,
    phone_number character varying(20) NOT NULL,
    client_id bigint,
    case_id bigint,
    message_sent boolean NOT NULL,
    full_name character varying(80)
);


--
-- Name: campaign_whatsapp_push_leads_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.campaign_whatsapp_push_leads_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: campaign_whatsapp_push_leads_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.campaign_whatsapp_push_leads_id_seq OWNED BY public.campaign_whatsapp_push_leads.id;


--
-- Name: campaigns; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.campaigns (
    id integer NOT NULL,
    company_id integer,
    name text NOT NULL,
    start_date date,
    end_date date,
    description text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    is_active boolean NOT NULL,
    funnel_id integer
);


--
-- Name: campaigns_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.campaigns_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: campaigns_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.campaigns_id_seq OWNED BY public.campaigns.id;


--
-- Name: cantons; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cantons (
    id integer NOT NULL,
    province_code text NOT NULL,
    country_code character(2) NOT NULL,
    code text NOT NULL,
    name text NOT NULL
);


--
-- Name: cantons_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.cantons_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cantons_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.cantons_id_seq OWNED BY public.cantons.id;


--
-- Name: case_funnel; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.case_funnel (
    id bigint NOT NULL,
    case_id integer NOT NULL,
    funnel_id integer,
    from_stage_id integer,
    to_stage_id integer,
    note text,
    changed_by integer NOT NULL,
    changed_at timestamp without time zone DEFAULT now() NOT NULL,
    action text DEFAULT 'move'::text NOT NULL,
    CONSTRAINT chk_cf_from_to_diff CHECK ((((action = 'assign'::text) AND (from_stage_id IS NULL) AND (to_stage_id IS NULL)) OR ((action = 'move'::text) AND (from_stage_id IS NOT NULL) AND (to_stage_id IS NOT NULL) AND (from_stage_id <> to_stage_id)) OR ((action = 'move'::text) AND (from_stage_id IS NULL) AND (to_stage_id IS NOT NULL)) OR ((action = 'unassign'::text) AND (from_stage_id IS NOT NULL) AND (to_stage_id IS NULL)) OR ((action = 'close'::text) AND (((from_stage_id IS NULL) AND (to_stage_id IS NULL)) OR ((from_stage_id IS NOT NULL) AND (to_stage_id IS NULL)))) OR (action = 'note'::text)))
);


--
-- Name: case_funnel_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.case_funnel_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: case_funnel_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.case_funnel_id_seq OWNED BY public.case_funnel.id;


--
-- Name: case_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.case_items (
    id integer NOT NULL,
    case_id integer NOT NULL,
    item_id integer NOT NULL,
    price numeric(14,2) DEFAULT 0 NOT NULL,
    quantity numeric(10,2) DEFAULT 1 NOT NULL,
    notes text,
    acquired boolean DEFAULT false,
    created_by integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: case_items_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.case_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: case_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.case_items_id_seq OWNED BY public.case_items.id;


--
-- Name: case_notes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.case_notes (
    id integer NOT NULL,
    case_id integer,
    author_id integer,
    note text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: case_notes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.case_notes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: case_notes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.case_notes_id_seq OWNED BY public.case_notes.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id integer NOT NULL,
    email text NOT NULL,
    full_name text,
    phone text,
    password_hash text,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    profile_image_url character varying,
    is_super_user boolean DEFAULT false
);


--
-- Name: case_notes_view; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.case_notes_view AS
 SELECT cn.id,
    cn.case_id,
    cn.author_id,
    u.full_name AS author_name,
    u.email AS author_email,
    cn.note,
    cn.created_at
   FROM (public.case_notes cn
     LEFT JOIN public.users u ON ((cn.author_id = u.id)));


--
-- Name: cases; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cases (
    id integer NOT NULL,
    client_id integer,
    campaign_id integer,
    company_id integer,
    department_id integer,
    agent_id integer,
    funnel_id integer,
    funnel_stage text,
    status text,
    started_at timestamp without time zone,
    closed_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    sender_id character varying(50),
    channel_id bigint,
    current_stage_id integer,
    channel_integration_id bigint,
    is_non_commercial boolean DEFAULT false,
    manual_starting_lead boolean DEFAULT false,
    payment_found boolean,
    CONSTRAINT cases_status_check CHECK ((status = ANY (ARRAY['open'::text, 'in_progress'::text, 'closed'::text, 'cancelled'::text])))
);


--
-- Name: cases_backup_before_sender_fix; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cases_backup_before_sender_fix (
    id integer,
    client_id integer,
    campaign_id integer,
    company_id integer,
    department_id integer,
    agent_id integer,
    funnel_id integer,
    funnel_stage text,
    status text,
    started_at timestamp without time zone,
    closed_at timestamp without time zone,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    sender_id character varying(50),
    channel_id bigint,
    current_stage_id integer,
    channel_integration_id bigint,
    is_non_commercial boolean,
    manual_starting_lead boolean,
    payment_found boolean
);


--
-- Name: cases_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.cases_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cases_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.cases_id_seq OWNED BY public.cases.id;


--
-- Name: channel_agent_client; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.channel_agent_client (
    id bigint NOT NULL,
    channel_id bigint NOT NULL,
    agent_id bigint NOT NULL,
    client_id bigint NOT NULL,
    department_id bigint NOT NULL
);


--
-- Name: channel_agent_client_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.channel_agent_client_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: channel_agent_client_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.channel_agent_client_id_seq OWNED BY public.channel_agent_client.id;


--
-- Name: channel_integrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.channel_integrations (
    id integer NOT NULL,
    company_id integer,
    channel_id integer,
    webhook_url text NOT NULL,
    access_token text,
    app_identifier text,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    is_non_commercial boolean DEFAULT false,
    integration_name text,
    department_id bigint,
    analyze_incoming_images boolean DEFAULT false
);


--
-- Name: channel_integrations_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.channel_integrations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: channel_integrations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.channel_integrations_id_seq OWNED BY public.channel_integrations.id;


--
-- Name: channel_whatsapp_template; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.channel_whatsapp_template (
    id bigint NOT NULL,
    department_id bigint NOT NULL,
    template_name character varying(255) NOT NULL,
    language character varying(10) NOT NULL,
    active boolean DEFAULT true NOT NULL,
    template_url_webhook character varying(500),
    channel_integration_id bigint,
    template_description text
);


--
-- Name: channel_whatsapp_template_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.channel_whatsapp_template_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: channel_whatsapp_template_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.channel_whatsapp_template_id_seq OWNED BY public.channel_whatsapp_template.id;


--
-- Name: channels; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.channels (
    id integer NOT NULL,
    code text NOT NULL,
    name text NOT NULL,
    description text
);


--
-- Name: channels_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.channels_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: channels_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.channels_id_seq OWNED BY public.channels.id;


--
-- Name: client_campaigns; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.client_campaigns (
    id integer NOT NULL,
    client_id integer,
    campaign_id integer,
    funnel_id integer,
    stage_detail text,
    assigned_agent_id integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: client_campaigns_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.client_campaigns_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: client_campaigns_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.client_campaigns_id_seq OWNED BY public.client_campaigns.id;


--
-- Name: client_files; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.client_files (
    id integer NOT NULL,
    client_id integer,
    uploader_id integer,
    file_url text,
    description text,
    mime_type text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: client_files_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.client_files_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: client_files_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.client_files_id_seq OWNED BY public.client_files.id;


--
-- Name: client_notes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.client_notes (
    id integer NOT NULL,
    client_id integer,
    author_id integer,
    note text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: client_notes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.client_notes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: client_notes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.client_notes_id_seq OWNED BY public.client_notes.id;


--
-- Name: client_social_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.client_social_accounts (
    id integer NOT NULL,
    client_id integer,
    channel_id integer,
    external_id text NOT NULL,
    username text,
    profile_pic text,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: client_social_accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.client_social_accounts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: client_social_accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.client_social_accounts_id_seq OWNED BY public.client_social_accounts.id;


--
-- Name: clients; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.clients (
    id integer NOT NULL,
    external_id text,
    full_name text,
    email text,
    phone text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    country_id integer,
    province_id integer,
    canton_id integer,
    district_id integer,
    address_detail text,
    postal_code text,
    is_citizen boolean DEFAULT false
);


--
-- Name: clients_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.clients_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: clients_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.clients_id_seq OWNED BY public.clients.id;


--
-- Name: companies; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.companies (
    id integer NOT NULL,
    name text NOT NULL,
    industry text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: companies_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.companies_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: companies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.companies_id_seq OWNED BY public.companies.id;


--
-- Name: countries; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.countries (
    id integer NOT NULL,
    iso_code character(2) NOT NULL,
    name text NOT NULL,
    phone_code text,
    currency_code character(3)
);


--
-- Name: countries_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.countries_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: countries_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.countries_id_seq OWNED BY public.countries.id;


--
-- Name: custom_field_definitions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.custom_field_definitions (
    id integer NOT NULL,
    entity_name text NOT NULL,
    field_key text NOT NULL,
    label text NOT NULL,
    field_type text NOT NULL,
    is_required boolean DEFAULT false,
    is_active boolean DEFAULT true,
    sort_order integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT custom_field_definitions_field_type_check CHECK ((field_type = ANY (ARRAY['text'::text, 'integer'::text, 'decimal'::text, 'date'::text, 'boolean'::text, 'collection'::text])))
);


--
-- Name: custom_field_definitions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.custom_field_definitions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: custom_field_definitions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.custom_field_definitions_id_seq OWNED BY public.custom_field_definitions.id;


--
-- Name: custom_field_values; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.custom_field_values (
    id integer NOT NULL,
    field_id integer NOT NULL,
    entity_name text NOT NULL,
    entity_id integer NOT NULL,
    value_text text,
    value_integer integer,
    value_decimal numeric(18,2),
    value_boolean boolean,
    value_date date,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: custom_field_values_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.custom_field_values_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: custom_field_values_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.custom_field_values_id_seq OWNED BY public.custom_field_values.id;


--
-- Name: custom_list_definitions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.custom_list_definitions (
    id bigint NOT NULL,
    list_name text NOT NULL,
    code_label text NOT NULL,
    description_label text NOT NULL,
    entity_name text NOT NULL,
    list_label text
);


--
-- Name: custom_list_definitions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.custom_list_definitions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: custom_list_definitions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.custom_list_definitions_id_seq OWNED BY public.custom_list_definitions.id;


--
-- Name: custom_list_entity_value; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.custom_list_entity_value (
    id bigint NOT NULL,
    entity_name text NOT NULL,
    entity_id bigint NOT NULL,
    list_value bigint NOT NULL,
    list_id bigint
);


--
-- Name: custom_list_entity_value_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.custom_list_entity_value_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: custom_list_entity_value_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.custom_list_entity_value_id_seq OWNED BY public.custom_list_entity_value.id;


--
-- Name: custom_list_values; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.custom_list_values (
    id bigint NOT NULL,
    list_id bigint NOT NULL,
    code_value text,
    description_value text
);


--
-- Name: custom_list_values_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.custom_list_values_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: custom_list_values_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.custom_list_values_id_seq OWNED BY public.custom_list_values.id;


--
-- Name: departments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.departments (
    id integer NOT NULL,
    company_id integer,
    name text NOT NULL,
    description text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: departments_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.departments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: departments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.departments_id_seq OWNED BY public.departments.id;


--
-- Name: districts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.districts (
    id integer NOT NULL,
    country_code character(2) NOT NULL,
    canton_code text NOT NULL,
    code text,
    name text NOT NULL,
    latitude numeric(10,7),
    longitude numeric(10,7),
    postal_code text
);


--
-- Name: districts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.districts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: districts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.districts_id_seq OWNED BY public.districts.id;


--
-- Name: erp_payment_confirmation; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.erp_payment_confirmation (
    id_record bigint NOT NULL,
    id integer NOT NULL,
    status text NOT NULL,
    harmony_state text NOT NULL,
    bank_name text NOT NULL,
    reference_number text NOT NULL,
    payment_date date NOT NULL,
    payment_time time without time zone NOT NULL,
    amount numeric(15,2) NOT NULL,
    client_document text NOT NULL,
    client_name text,
    contract_id integer,
    contract_payment_method_id integer,
    contract_payment_method_name text,
    raw_text text,
    receipt_type text NOT NULL,
    receipt_base64 text,
    receipt_message text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    notificacion_pago_id bigint,
    custom_message text
);


--
-- Name: erp_payment_confirmation_id_record_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.erp_payment_confirmation_id_record_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: erp_payment_confirmation_id_record_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.erp_payment_confirmation_id_record_seq OWNED BY public.erp_payment_confirmation.id_record;


--
-- Name: funnel_stages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.funnel_stages (
    id integer NOT NULL,
    funnel_id integer NOT NULL,
    name text NOT NULL,
    code text,
    "position" integer NOT NULL,
    color_hex text,
    is_terminal boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    is_won boolean DEFAULT false NOT NULL,
    is_lost boolean DEFAULT false NOT NULL
);


--
-- Name: funnel_stages_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.funnel_stages_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: funnel_stages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.funnel_stages_id_seq OWNED BY public.funnel_stages.id;


--
-- Name: funnels; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.funnels (
    id integer NOT NULL,
    name text NOT NULL,
    stage_order integer,
    description text,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--
-- Name: funnels_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.funnels_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: funnels_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.funnels_id_seq OWNED BY public.funnels.id;


--
-- Name: integration_templates; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.integration_templates (
    id integer NOT NULL,
    integration_id integer NOT NULL,
    template_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: integration_templates_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.integration_templates_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: integration_templates_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.integration_templates_id_seq OWNED BY public.integration_templates.id;


--
-- Name: item_departments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.item_departments (
    item_id integer NOT NULL,
    department_id integer NOT NULL
);


--
-- Name: items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.items (
    id integer NOT NULL,
    company_id integer,
    name text NOT NULL,
    description text,
    type text,
    item_price numeric(10,3) DEFAULT 0 NOT NULL,
    CONSTRAINT items_type_check CHECK ((type = ANY (ARRAY['product'::text, 'service'::text])))
);


--
-- Name: items_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.items_id_seq OWNED BY public.items.id;


--
-- Name: message_case_full_backup_before_sender_fix; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.message_case_full_backup_before_sender_fix (
    message_id integer,
    original_case_id integer,
    backed_up_at timestamp with time zone
);


--
-- Name: message_status; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.message_status (
    id bigint NOT NULL,
    channel_message_id text NOT NULL,
    message_status text NOT NULL,
    applied boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


--
-- Name: message_status_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.message_status_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: message_status_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.message_status_id_seq OWNED BY public.message_status.id;


--
-- Name: message_templates; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.message_templates (
    id integer NOT NULL,
    channel_id integer NOT NULL,
    template_name text NOT NULL,
    language_code character varying(10) NOT NULL,
    description text,
    category text,
    is_active boolean DEFAULT true,
    is_conversation_starter boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: message_templates_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.message_templates_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: message_templates_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.message_templates_id_seq OWNED BY public.message_templates.id;


--
-- Name: messages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.messages (
    id integer NOT NULL,
    case_id integer,
    sender_type text,
    message_type text,
    text_content text,
    file_url text,
    mime_type text,
    channel_message_id text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    base64_content text,
    agent_id bigint,
    has_error boolean,
    message_error text,
    message_read boolean,
    status text,
    CONSTRAINT chk_messages_text_content_required CHECK (((message_type <> 'text'::text) OR ((text_content IS NOT NULL) AND (btrim(text_content) <> ''::text)))),
    CONSTRAINT messages_message_type_check CHECK ((message_type = ANY (ARRAY['text'::text, 'image'::text, 'file'::text, 'audio'::text, 'sticker'::text]))),
    CONSTRAINT messages_sender_type_check CHECK ((sender_type = ANY (ARRAY['client'::text, 'agent'::text])))
);


--
-- Name: messages_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.messages_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: messages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.messages_id_seq OWNED BY public.messages.id;


--
-- Name: mv_cases_with_channels; Type: MATERIALIZED VIEW; Schema: public; Owner: -
--

CREATE MATERIALIZED VIEW public.mv_cases_with_channels AS
 SELECT c.id AS case_id,
    c.client_id,
    cl.full_name AS client_name,
    c.campaign_id,
    c.company_id,
    c.department_id,
    c.agent_id,
    (c.agent_id IS NOT NULL) AS agent_assigned,
    u.full_name AS agent_full_name,
    c.funnel_id,
    (NULLIF(c.funnel_stage, ''::text))::integer AS funnel_stage,
    c.status,
    c.channel_id,
    ch.code AS channel_code,
    ch.name AS channel_name,
    ci.integration_name,
    c.started_at,
    c.closed_at,
    c.created_at,
    c.updated_at,
    c.sender_id,
    lm.id AS last_message_id,
    lm.sender_type AS last_message_sender_type,
    lm.message_type AS last_message_type,
    lm.text_content AS last_message_text,
    lm.created_at AS last_message_at,
        CASE
            WHEN (lm.message_type = 'text'::text) THEN lm.text_content
            WHEN (lm.message_type = 'image'::text) THEN '[Imagen]'::text
            WHEN (lm.message_type = 'audio'::text) THEN '[Audio]'::text
            WHEN (lm.message_type = 'file'::text) THEN '[Archivo]'::text
            ELSE NULL::text
        END AS last_message_preview,
    COALESCE(msg_counts.client_messages, (0)::bigint) AS client_messages,
    COALESCE(unread.unread_count, (0)::bigint) AS unread_count
   FROM (((((((public.cases c
     LEFT JOIN public.channels ch ON ((ch.id = c.channel_id)))
     LEFT JOIN public.channel_integrations ci ON ((ci.id = c.channel_integration_id)))
     LEFT JOIN public.clients cl ON ((cl.id = c.client_id)))
     LEFT JOIN public.users u ON ((u.id = c.agent_id)))
     LEFT JOIN LATERAL ( SELECT m.id,
            m.case_id,
            m.sender_type,
            m.message_type,
            m.text_content,
            m.created_at
           FROM public.messages m
          WHERE (m.case_id = c.id)
          ORDER BY m.id DESC
         LIMIT 1) lm ON (true))
     LEFT JOIN ( SELECT messages.case_id,
            count(*) AS client_messages
           FROM public.messages
          WHERE (messages.sender_type = 'client'::text)
          GROUP BY messages.case_id) msg_counts ON ((msg_counts.case_id = c.id)))
     LEFT JOIN ( SELECT messages.case_id,
            count(*) AS unread_count
           FROM public.messages
          WHERE (messages.message_read = false)
          GROUP BY messages.case_id) unread ON ((unread.case_id = c.id)))
  WITH NO DATA;


--
-- Name: permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.permissions (
    id integer NOT NULL,
    code text NOT NULL,
    description text
);


--
-- Name: permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.permissions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.permissions_id_seq OWNED BY public.permissions.id;


--
-- Name: provinces; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.provinces (
    id integer NOT NULL,
    country_code character(2) NOT NULL,
    code text NOT NULL,
    name text NOT NULL
);


--
-- Name: provinces_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.provinces_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: provinces_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.provinces_id_seq OWNED BY public.provinces.id;


--
-- Name: qr_leads; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.qr_leads (
    id integer NOT NULL,
    company_id integer NOT NULL,
    campaign_id integer NOT NULL,
    department_id integer,
    user_id integer,
    client_id integer,
    contact_phone character varying(50) NOT NULL,
    status text DEFAULT 'pending'::text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT qr_leads_status_check CHECK ((status = ANY (ARRAY['pending'::text, 'contacted'::text, 'converted'::text, 'discarded'::text])))
);


--
-- Name: qr_leads_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.qr_leads_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: qr_leads_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.qr_leads_id_seq OWNED BY public.qr_leads.id;


--
-- Name: receipt_process; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.receipt_process (
    id_record bigint NOT NULL,
    process_date timestamp with time zone DEFAULT now() NOT NULL,
    records_processed integer NOT NULL,
    total_time_ms bigint NOT NULL,
    avg_time_per_record_ms numeric(10,2) NOT NULL
);


--
-- Name: receipt_process_id_record_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.receipt_process_id_record_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: receipt_process_id_record_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.receipt_process_id_record_seq OWNED BY public.receipt_process.id_record;


--
-- Name: receipt_results; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.receipt_results (
    id integer NOT NULL,
    case_id bigint,
    status character varying(30) DEFAULT 'new'::character varying,
    bank_name character varying(100),
    transaction_type character varying(100),
    reference_number character varying(100),
    date character varying(100),
    "time" character varying(20),
    amount double precision,
    amount_sent double precision,
    sender_name character varying(200),
    sender_phone character varying(50),
    receiver_name character varying(200),
    receiver_phone character varying(50),
    origin_account character varying(100),
    destination_account character varying(100),
    description text,
    raw_text text,
    warnings text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    message_id bigint
);


--
-- Name: receipt_results_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.receipt_results_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: receipt_results_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.receipt_results_id_seq OWNED BY public.receipt_results.id;


--
-- Name: role_permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.role_permissions (
    role_id integer NOT NULL,
    permission_id integer NOT NULL
);


--
-- Name: roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.roles (
    id integer NOT NULL,
    name text NOT NULL,
    description text,
    is_agent boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.roles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- Name: seats_sale; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.seats_sale (
    id integer NOT NULL,
    seat_number text NOT NULL,
    zone text NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


--
-- Name: seats_sale_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.seats_sale_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: seats_sale_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.seats_sale_id_seq OWNED BY public.seats_sale.id;


--
-- Name: settings; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.settings (
    id bigint NOT NULL,
    value_code character varying(200) NOT NULL,
    text_value text,
    integer_value bigint,
    number_value numeric(18,3),
    bool_value boolean,
    is_active boolean NOT NULL
);


--
-- Name: settings_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.settings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: settings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.settings_id_seq OWNED BY public.settings.id;


--
-- Name: user_channel_permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_channel_permissions (
    id integer NOT NULL,
    user_id integer,
    company_id integer,
    channel_id integer,
    can_receive boolean DEFAULT true,
    can_send boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: user_channel_permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_channel_permissions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_channel_permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.user_channel_permissions_id_seq OWNED BY public.user_channel_permissions.id;


--
-- Name: user_company_roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_company_roles (
    id integer NOT NULL,
    user_id integer,
    company_id integer,
    role_id integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: user_company_roles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_company_roles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_company_roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.user_company_roles_id_seq OWNED BY public.user_company_roles.id;


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: vw_agent_department_assignments; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_agent_department_assignments AS
 SELECT d.company_id,
    d.id AS department_id,
    d.name AS department_name,
    a.user_id,
        CASE
            WHEN (ada.id IS NOT NULL) THEN true
            ELSE false
        END AS department_assigned
   FROM ((public.departments d
     CROSS JOIN public.agents a)
     LEFT JOIN public.agent_department_assignments ada ON (((ada.department_id = d.id) AND (ada.agent_id = a.user_id))));


--
-- Name: vw_agent_users; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_agent_users AS
 SELECT u.id,
    u.email,
    u.full_name,
    u.phone,
    u.is_active,
    u.profile_image_url
   FROM (public.users u
     JOIN public.agents a ON ((a.user_id = u.id)));


--
-- Name: vw_agent_department_information; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_agent_department_information AS
 SELECT au.id AS agent_id,
    au.full_name AS agent_name,
    ada.company_id,
    ada.department_name,
    ada.department_id,
    au.profile_image_url
   FROM (public.vw_agent_department_assignments ada
     JOIN public.vw_agent_users au ON ((ada.user_id = au.id)))
  WHERE ((ada.department_assigned = true) AND (au.is_active = true));


--
-- Name: vw_campaign_funnel_summary; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_campaign_funnel_summary AS
 SELECT ca.id AS campaign_id,
    ca.name AS campaign_name,
    ca.company_id,
    f.id AS funnel_id,
    f.name AS funnel_name,
    fs.id AS stage_id,
    fs.name AS stage_name,
    fs.code AS stage_code,
    fs."position",
    fs.color_hex,
    fs.is_won,
    fs.is_lost,
    fs.is_terminal,
    count(cs.id) AS total_cases
   FROM (((public.campaigns ca
     JOIN public.funnels f ON ((f.id = ca.funnel_id)))
     JOIN public.funnel_stages fs ON ((fs.funnel_id = f.id)))
     LEFT JOIN public.cases cs ON (((cs.campaign_id = ca.id) AND (cs.current_stage_id = fs.id))))
  GROUP BY ca.id, ca.name, ca.company_id, f.id, f.name, fs.id, fs.name, fs.code, fs."position", fs.color_hex, fs.is_won, fs.is_lost, fs.is_terminal
  ORDER BY ca.id, fs."position";


--
-- Name: vw_campaigns_with_funnel; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_campaigns_with_funnel AS
 SELECT c.id AS campaign_id,
    c.company_id,
    c.name AS campaign_name,
    c.start_date,
    c.end_date,
    c.description,
    c.created_at,
    c.is_active,
    c.funnel_id,
    f.name AS funnel_name
   FROM (public.campaigns c
     LEFT JOIN public.funnels f ON ((f.id = c.funnel_id)));


--
-- Name: vw_case_channel_integration; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_case_channel_integration AS
 SELECT c.id AS case_id,
    c.company_id,
    c.channel_id,
    c.sender_id,
    c.status,
    c.started_at,
    c.closed_at,
    c.channel_integration_id,
    ch.name AS channel_name,
    ch.code AS channel_code,
    ci.webhook_url,
    ci.access_token,
    ci.app_identifier,
    ci.is_active AS integration_is_active,
    ci.updated_at AS integration_updated_at
   FROM ((public.cases c
     JOIN public.channel_integrations ci ON ((ci.id = c.channel_integration_id)))
     JOIN public.channels ch ON ((ch.id = ci.channel_id)));


--
-- Name: vw_case_current_stage; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_case_current_stage AS
 SELECT DISTINCT ON (cf.case_id) cf.case_id,
    cf.funnel_id,
    f.name AS funnel_name,
    cf.to_stage_id AS current_stage_id,
        CASE
            WHEN ((cf.to_stage_id IS NULL) AND (cf.action = 'assign'::text)) THEN 'NA'::text
            ELSE fst.name
        END AS current_stage_name,
    cf.changed_at AS last_changed_at,
    cf.changed_by AS last_changed_by,
    u.full_name AS last_changed_by_label,
    cf.action,
    fst.color_hex
   FROM (((public.case_funnel cf
     LEFT JOIN public.funnels f ON ((f.id = cf.funnel_id)))
     LEFT JOIN public.funnel_stages fst ON ((fst.id = cf.to_stage_id)))
     LEFT JOIN public.users u ON ((u.id = cf.changed_by)))
  ORDER BY cf.case_id, cf.changed_at DESC;


--
-- Name: vw_case_dashboard_by_company; Type: MATERIALIZED VIEW; Schema: public; Owner: -
--

CREATE MATERIALIZED VIEW public.vw_case_dashboard_by_company AS
 WITH channel_stats AS (
         SELECT cs_1.company_id,
            ch.id AS channel_id,
            ch.name AS channel_name,
            count(cs_1.id) FILTER (WHERE (cs_1.status = 'open'::text)) AS open_cases,
            count(cs_1.id) FILTER (WHERE (cs_1.status = 'closed'::text)) AS closed_cases
           FROM (public.cases cs_1
             LEFT JOIN public.channels ch ON ((ch.id = cs_1.channel_id)))
          GROUP BY cs_1.company_id, ch.id, ch.name
        ), agent_stats AS (
         SELECT cs_1.company_id,
            u.id AS agent_id,
            u.full_name AS agent_name,
            count(cs_1.id) FILTER (WHERE (cs_1.status = 'open'::text)) AS open_cases,
            count(cs_1.id) FILTER (WHERE (cs_1.status = 'closed'::text)) AS closed_cases,
            round(avg((EXTRACT(epoch FROM (cs_1.closed_at - cs_1.started_at)) / (3600)::numeric)), 2) AS avg_close_hours
           FROM (public.cases cs_1
             LEFT JOIN public.users u ON ((u.id = cs_1.agent_id)))
          WHERE (cs_1.agent_id IS NOT NULL)
          GROUP BY cs_1.company_id, u.id, u.full_name
        ), oldest_open AS (
         SELECT cs_1.company_id,
            cs_1.id AS case_id,
            cl.full_name AS client_name,
            cs_1.started_at
           FROM (public.cases cs_1
             LEFT JOIN public.clients cl ON ((cl.id = cs_1.client_id)))
          WHERE (cs_1.status = 'open'::text)
        )
 SELECT c.id AS company_id,
    c.name AS company_name,
    count(cs.id) AS total_cases,
    count(cs.id) FILTER (WHERE (cs.status = 'open'::text)) AS open_cases,
    count(cs.id) FILTER (WHERE (cs.status = 'closed'::text)) AS closed_cases,
    count(cs.id) FILTER (WHERE ((cs.status = 'closed'::text) AND (date(cs.closed_at) = CURRENT_DATE))) AS closed_today,
    count(cs.id) FILTER (WHERE ((cs.status = 'open'::text) AND (date(cs.created_at) = CURRENT_DATE))) AS opened_today,
    count(cs.id) FILTER (WHERE (cs.status = 'cancelled'::text)) AS cancelled_cases,
    count(cs.id) FILTER (WHERE (cs.agent_id IS NULL)) AS unassigned_agents,
    count(cs.id) FILTER (WHERE (cs.client_id IS NULL)) AS unassigned_clients,
    round(avg((EXTRACT(epoch FROM (cs.closed_at - cs.started_at)) / (3600)::numeric)), 2) AS avg_close_hours,
    ( SELECT json_agg(json_build_object('channel_id', ch.channel_id, 'channel_name', ch.channel_name, 'open_cases', ch.open_cases, 'closed_cases', ch.closed_cases)) AS json_agg
           FROM channel_stats ch
          WHERE (ch.company_id = c.id)) AS cases_by_channel,
    ( SELECT json_agg(json_build_object('agent_id', a.agent_id, 'agent_name', a.agent_name, 'open_cases', a.open_cases, 'closed_cases', a.closed_cases, 'avg_close_hours', a.avg_close_hours)) AS json_agg
           FROM agent_stats a
          WHERE (a.company_id = c.id)) AS cases_by_agent,
    ( SELECT json_agg(json_build_object('case_id', o.case_id, 'client_name', o.client_name, 'started_at', o.started_at) ORDER BY o.started_at) AS json_agg
           FROM oldest_open o
          WHERE (o.company_id = c.id)
         LIMIT 5) AS oldest_open_cases
   FROM (public.companies c
     LEFT JOIN public.cases cs ON ((cs.company_id = c.id)))
  GROUP BY c.id, c.name
  ORDER BY c.id
  WITH NO DATA;


--
-- Name: vw_case_dashboard_by_company_with_department; Type: MATERIALIZED VIEW; Schema: public; Owner: -
--

CREATE MATERIALIZED VIEW public.vw_case_dashboard_by_company_with_department AS
 WITH channel_stats AS (
         SELECT cs_1.company_id,
            cs_1.department_id,
            ch.id AS channel_id,
            ch.name AS channel_name,
            count(cs_1.id) FILTER (WHERE (cs_1.status = 'open'::text)) AS open_cases,
            count(cs_1.id) FILTER (WHERE (cs_1.status = 'closed'::text)) AS closed_cases
           FROM (public.cases cs_1
             LEFT JOIN public.channels ch ON ((ch.id = cs_1.channel_id)))
          GROUP BY cs_1.company_id, cs_1.department_id, ch.id, ch.name
        ), agent_stats AS (
         SELECT cs_1.company_id,
            cs_1.department_id,
            u.id AS agent_id,
            u.full_name AS agent_name,
            count(cs_1.id) FILTER (WHERE (cs_1.status = 'open'::text)) AS open_cases,
            count(cs_1.id) FILTER (WHERE (cs_1.status = 'closed'::text)) AS closed_cases,
            round(avg((EXTRACT(epoch FROM (cs_1.closed_at - cs_1.started_at)) / 3600.0)), 2) AS avg_close_hours
           FROM (public.cases cs_1
             LEFT JOIN public.users u ON ((u.id = cs_1.agent_id)))
          WHERE (cs_1.agent_id IS NOT NULL)
          GROUP BY cs_1.company_id, cs_1.department_id, u.id, u.full_name
        ), oldest_open AS (
         SELECT cs_1.company_id,
            cs_1.department_id,
            cs_1.id AS case_id,
            cl.full_name AS client_name,
            cs_1.created_at
           FROM (public.cases cs_1
             LEFT JOIN public.clients cl ON ((cl.id = cs_1.client_id)))
          WHERE (cs_1.status = 'open'::text)
        )
 SELECT c.id AS company_id,
    c.name AS company_name,
    cs.department_id,
    count(cs.id) AS total_cases,
    count(cs.id) FILTER (WHERE (cs.status = 'open'::text)) AS open_cases,
    count(cs.id) FILTER (WHERE (cs.status = 'closed'::text)) AS closed_cases,
    count(cs.id) FILTER (WHERE ((cs.status = 'closed'::text) AND (date(cs.closed_at) = CURRENT_DATE))) AS closed_today,
    count(cs.id) FILTER (WHERE ((cs.status = 'open'::text) AND (date(cs.created_at) = CURRENT_DATE))) AS opened_today,
    count(cs.id) FILTER (WHERE (cs.status = 'cancelled'::text)) AS cancelled_cases,
    count(cs.id) FILTER (WHERE (cs.agent_id IS NULL)) AS unassigned_agents,
    count(cs.id) FILTER (WHERE (cs.client_id IS NULL)) AS unassigned_clients,
    round(avg((EXTRACT(epoch FROM (cs.closed_at - cs.started_at)) / 3600.0)), 2) AS avg_close_hours,
    ( SELECT json_agg(json_build_object('channel_id', ch.channel_id, 'channel_name', ch.channel_name, 'department_id', ch.department_id, 'open_cases', ch.open_cases, 'closed_cases', ch.closed_cases)) AS json_agg
           FROM channel_stats ch
          WHERE ((ch.company_id = c.id) AND (ch.department_id = cs.department_id))) AS cases_by_channel,
    ( SELECT json_agg(json_build_object('agent_id', a.agent_id, 'agent_name', a.agent_name, 'department_id', a.department_id, 'open_cases', a.open_cases, 'closed_cases', a.closed_cases, 'avg_close_hours', a.avg_close_hours)) AS json_agg
           FROM agent_stats a
          WHERE ((a.company_id = c.id) AND (a.department_id = cs.department_id))) AS cases_by_agent,
    ( SELECT json_agg(json_build_object('case_id', o.case_id, 'client_name', o.client_name, 'department_id', o.department_id, 'created_at', o.created_at) ORDER BY o.created_at) AS json_agg
           FROM oldest_open o
          WHERE ((o.company_id = c.id) AND (o.department_id = cs.department_id))
         LIMIT 5) AS oldest_open_cases
   FROM (public.companies c
     LEFT JOIN public.cases cs ON ((cs.company_id = c.id)))
  GROUP BY c.id, c.name, cs.department_id
  ORDER BY c.id
  WITH NO DATA;


--
-- Name: vw_case_general_information; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_case_general_information AS
 SELECT c.id AS case_id,
    c.company_id,
    c.client_id,
    c.department_id,
    c.campaign_id,
    c.agent_id,
    c.funnel_id,
    c.status,
    c.channel_id,
    c.sender_id,
    cs.current_stage_id,
    cs.current_stage_name,
    cs.last_changed_by_label,
    cs.action,
    cs.color_hex,
    cl.full_name AS client_name,
    cl.email,
    cn.name AS campaign_name,
    ch.name AS channel_name,
    ch.code AS channel_code,
    d.name AS department_name,
    u.full_name AS agent_name,
    c.channel_integration_id,
    c.manual_starting_lead,
    count(m.id) AS client_messages
   FROM (((((((public.cases c
     LEFT JOIN public.vw_case_current_stage cs ON ((c.id = cs.case_id)))
     LEFT JOIN public.clients cl ON ((c.client_id = cl.id)))
     LEFT JOIN public.campaigns cn ON ((c.campaign_id = cn.id)))
     LEFT JOIN public.channels ch ON ((c.channel_id = ch.id)))
     LEFT JOIN public.departments d ON ((c.department_id = d.id)))
     LEFT JOIN public.users u ON ((c.agent_id = u.id)))
     LEFT JOIN public.messages m ON (((m.case_id = c.id) AND (m.sender_type = 'client'::text))))
  GROUP BY c.id, c.company_id, c.client_id, c.department_id, c.campaign_id, c.agent_id, c.funnel_id, c.status, c.channel_id, c.sender_id, cs.current_stage_id, cs.current_stage_name, cs.last_changed_by_label, cs.action, cs.color_hex, cl.full_name, cl.email, cn.name, ch.name, ch.code, d.name, u.full_name, c.channel_integration_id, c.manual_starting_lead
  ORDER BY c.id DESC;


--
-- Name: vw_case_items_detail; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_case_items_detail AS
 SELECT ci.id AS case_item_id,
    ci.case_id,
    ci.item_id,
    c.company_id,
    c.department_id,
    c.campaign_id,
    c.client_id,
    cl.full_name AS client_name,
    cl.email AS client_email,
    cl.phone AS client_phone,
    i.name AS item_name,
    i.description AS item_description,
    i.type AS item_type,
    ci.price,
    ci.quantity,
    (ci.price * ci.quantity) AS total_amount,
    ci.acquired,
    ci.notes,
    ci.created_by,
    u.full_name AS created_by_name,
    ci.created_at,
    c.status AS case_status,
    c.funnel_stage,
    c.started_at,
    c.closed_at
   FROM ((((public.case_items ci
     JOIN public.cases c ON ((ci.case_id = c.id)))
     JOIN public.items i ON ((ci.item_id = i.id)))
     LEFT JOIN public.users u ON ((ci.created_by = u.id)))
     LEFT JOIN public.clients cl ON ((c.client_id = cl.id)));


--
-- Name: vw_cases_images_without_receipt; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_cases_images_without_receipt AS
 SELECT case_id,
    base64_content,
    created_at
   FROM public.messages m
  WHERE ((message_type = 'image'::text) AND (sender_type = 'client'::text) AND (base64_content IS NOT NULL) AND (NOT (EXISTS ( SELECT 1
           FROM public.receipt_results rr
          WHERE (rr.case_id = m.case_id)))));


--
-- Name: vw_cases_with_channels; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_cases_with_channels AS
 SELECT c.id AS case_id,
    c.client_id,
    cl.full_name AS client_name,
    c.campaign_id,
    c.company_id,
    c.department_id,
    c.agent_id,
        CASE
            WHEN (c.agent_id IS NULL) THEN false
            ELSE true
        END AS agent_assigned,
    u.full_name AS agent_full_name,
    c.funnel_id,
    c.funnel_stage,
    c.status,
    c.channel_id,
    ch.code AS channel_code,
    ch.name AS channel_name,
    ch.description AS channel_description,
    ci.integration_name,
    c.started_at,
    c.closed_at,
    c.created_at,
    c.updated_at,
    c.sender_id,
    c.payment_found,
    lm.id AS last_message_id,
    lm.sender_type AS last_message_sender_type,
    lm.message_type AS last_message_type,
    lm.text_content AS last_message_text,
    lm.file_url AS last_message_file_url,
    lm.mime_type AS last_message_mime_type,
    lm.created_at AS last_message_at,
        CASE
            WHEN (lm.message_type = 'text'::text) THEN lm.text_content
            WHEN (lm.message_type = 'image'::text) THEN '[Imagen]'::text
            WHEN (lm.message_type = 'audio'::text) THEN '[Audio]'::text
            WHEN (lm.message_type = 'file'::text) THEN '[Archivo]'::text
            ELSE NULL::text
        END AS last_message_preview,
        CASE
            WHEN (lm.message_type = ANY (ARRAY['image'::text, 'audio'::text, 'file'::text])) THEN true
            ELSE false
        END AS last_message_is_media,
    c.is_non_commercial,
    c.channel_integration_id,
    c.manual_starting_lead,
    COALESCE(msg_counts.client_messages, (0)::bigint) AS client_messages,
    COALESCE(unread.unread_count, (0)::bigint) AS unread_count
   FROM (((((((public.cases c
     LEFT JOIN public.channels ch ON ((ch.id = c.channel_id)))
     LEFT JOIN public.channel_integrations ci ON ((ci.id = c.channel_integration_id)))
     LEFT JOIN public.clients cl ON ((cl.id = c.client_id)))
     LEFT JOIN public.users u ON ((u.id = c.agent_id)))
     LEFT JOIN LATERAL ( SELECT msg.id,
            msg.case_id,
            msg.sender_type,
            msg.message_type,
            msg.text_content,
            msg.file_url,
            msg.mime_type,
            msg.created_at
           FROM public.messages msg
          WHERE (msg.case_id = c.id)
          ORDER BY msg.id DESC
         LIMIT 1) lm ON (true))
     LEFT JOIN ( SELECT messages.case_id,
            count(*) AS client_messages
           FROM public.messages
          WHERE (messages.sender_type = 'client'::text)
          GROUP BY messages.case_id) msg_counts ON ((msg_counts.case_id = c.id)))
     LEFT JOIN ( SELECT m.case_id,
            count(*) AS unread_count
           FROM public.messages m
          WHERE (m.message_read = false)
          GROUP BY m.case_id) unread ON ((unread.case_id = c.id)));


--
-- Name: vw_channel_integrations; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_channel_integrations AS
 SELECT ci.id AS channel_integration_id,
    ci.company_id,
    ci.channel_id,
    ci.webhook_url,
    ci.access_token,
    ci.app_identifier,
    ci.is_active,
    ci.created_at,
    ci.updated_at,
    ci.is_non_commercial,
    ci.integration_name,
    ci.department_id,
    d.name AS department_name,
    c.code AS channel_code,
    ci.analyze_incoming_images
   FROM ((public.channel_integrations ci
     JOIN public.channels c ON ((c.id = ci.channel_id)))
     LEFT JOIN public.departments d ON ((ci.department_id = d.id)));


--
-- Name: vw_channel_template_integrations; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_channel_template_integrations AS
 SELECT ch.id AS channel_id,
    ch.code AS channel_code,
    ch.name AS channel_name,
    mt.id AS template_id,
    mt.template_name,
    mt.description,
    mt.language_code,
    ci.id AS integration_id,
    ci.integration_name,
    ci.company_id,
    ci.department_id,
        CASE
            WHEN (it.id IS NOT NULL) THEN true
            ELSE false
        END AS is_linked
   FROM (((public.channels ch
     JOIN public.message_templates mt ON ((mt.channel_id = ch.id)))
     JOIN public.channel_integrations ci ON ((ci.channel_id = ch.id)))
     LEFT JOIN public.integration_templates it ON (((it.template_id = mt.id) AND (it.integration_id = ci.id))));


--
-- Name: vw_channel_whatsapp_templates; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_channel_whatsapp_templates AS
 SELECT DISTINCT t.id,
    t.department_id,
    t.template_name,
    t.language,
    t.active,
    t.template_url_webhook,
    ci.company_id,
    ci.channel_id
   FROM (public.channel_whatsapp_template t
     JOIN public.channel_integrations ci ON ((ci.department_id = t.department_id)))
  WHERE (t.active = true)
  ORDER BY t.template_name;


--
-- Name: vw_channels; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_channels AS
 SELECT ci.id AS integration_id,
    ci.company_id,
    c.id AS channel_id,
    c.code AS channel_code,
    c.name AS channel_name,
    c.description AS channel_description,
    ci.webhook_url,
    ci.access_token,
    ci.app_identifier,
    ci.is_active,
    ci.created_at,
    ci.updated_at,
    ci.is_non_commercial,
    ci.department_id
   FROM (public.channel_integrations ci
     JOIN public.channels c ON ((c.id = ci.channel_id)));


--
-- Name: vw_client_social_accounts; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_client_social_accounts AS
 SELECT c.id AS client_id,
    c.external_id AS client_external_id,
    c.full_name,
    c.email,
    c.phone,
    c.created_at AS client_created_at,
    c.updated_at AS client_updated_at,
    sa.id AS social_account_id,
    sa.channel_id,
    sa.external_id AS social_external_id,
    sa.username,
    sa.profile_pic,
    sa.is_active,
    sa.created_at AS social_created_at,
    sa.updated_at AS social_updated_at
   FROM (public.clients c
     LEFT JOIN public.client_social_accounts sa ON ((sa.client_id = c.id)));


--
-- Name: vw_company_channel_templates; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_company_channel_templates AS
 SELECT ci.company_id,
    ci.id AS channel_integration_id,
    ci.webhook_url,
    ci.access_token,
    ci.app_identifier,
    ci.is_active AS integration_active,
    ci.created_at AS integration_created_at,
    ci.updated_at AS integration_updated_at,
    ch.id AS channel_id,
    ch.code AS channel_code,
    ch.name AS channel_name,
    ch.description AS channel_description,
    cwt.id AS template_id,
    cwt.template_name,
    cwt.language,
    cwt.template_url_webhook,
    cwt.active AS template_active
   FROM ((public.channel_integrations ci
     JOIN public.channels ch ON ((ch.id = ci.channel_id)))
     JOIN public.channel_whatsapp_template cwt ON ((cwt.department_id = ci.department_id)));


--
-- Name: vw_company_users; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_company_users AS
 SELECT DISTINCT c.id AS company_id,
    c.name AS company_name,
    u.id AS user_id,
    u.full_name,
    u.email
   FROM ((public.user_company_roles ucr
     JOIN public.companies c ON ((ucr.company_id = c.id)))
     JOIN public.users u ON ((ucr.user_id = u.id)));


--
-- Name: vw_custom_fields; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_custom_fields AS
 SELECT d.id AS field_id,
    d.entity_name,
    d.field_key,
    d.label,
    d.field_type,
    d.is_required,
    d.is_active,
    d.sort_order,
    v.entity_id,
    v.value_text,
    v.value_integer,
    v.value_decimal,
    v.value_boolean,
    v.value_date
   FROM (public.custom_field_definitions d
     LEFT JOIN public.custom_field_values v ON (((v.field_id = d.id) AND (v.entity_name = d.entity_name))));


--
-- Name: vw_dashboard_campaigns_per_company; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_dashboard_campaigns_per_company AS
 SELECT cu.company_id,
    cu.company_name,
    c.id AS campaign_id,
    c.name AS campaign_name,
    c.is_active,
    count(DISTINCT cs.id) AS total_cases,
    count(DISTINCT cs.id) FILTER (WHERE (cs.status = 'open'::text)) AS open_cases,
    count(DISTINCT cs.id) FILTER (WHERE (cs.status = 'closed'::text)) AS closed_cases,
    count(DISTINCT cs.id) FILTER (WHERE (EXISTS ( SELECT 1
           FROM public.funnel_stages fs
          WHERE ((fs.id = cs.current_stage_id) AND (fs.is_won = true))))) AS won_cases,
    count(DISTINCT cs.id) FILTER (WHERE (EXISTS ( SELECT 1
           FROM public.funnel_stages fs
          WHERE ((fs.id = cs.current_stage_id) AND (fs.is_lost = true))))) AS lost_cases,
    round(((100.0 * (count(DISTINCT cs.id) FILTER (WHERE (EXISTS ( SELECT 1
           FROM public.funnel_stages fs
          WHERE ((fs.id = cs.current_stage_id) AND (fs.is_won = true))))))::numeric) / (NULLIF(count(DISTINCT cs.id), 0))::numeric), 2) AS conversion_rate
   FROM ((public.vw_company_users cu
     JOIN public.campaigns c ON ((c.company_id = cu.company_id)))
     LEFT JOIN public.cases cs ON ((cs.campaign_id = c.id)))
  GROUP BY cu.company_id, cu.company_name, c.id, c.name, c.is_active
  ORDER BY cu.company_id, c.id;


--
-- Name: vw_dashboard_general_by_company; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_dashboard_general_by_company AS
 SELECT c.id AS company_id,
    c.name AS company_name,
    count(DISTINCT ca.id) FILTER (WHERE (ca.is_active = true)) AS total_active_campaigns,
    count(DISTINCT cs.id) AS total_cases,
    count(DISTINCT cs.id) FILTER (WHERE (cs.status = 'closed'::text)) AS closed_cases,
    count(DISTINCT cs.id) FILTER (WHERE (EXISTS ( SELECT 1
           FROM public.funnel_stages fs
          WHERE ((fs.id = cs.current_stage_id) AND (fs.is_won = true))))) AS won_cases,
    round(((100.0 * (count(DISTINCT cs.id) FILTER (WHERE (EXISTS ( SELECT 1
           FROM public.funnel_stages fs
          WHERE ((fs.id = cs.current_stage_id) AND (fs.is_won = true))))))::numeric) / (NULLIF(count(DISTINCT cs.id), 0))::numeric), 2) AS conversion_rate,
    count(DISTINCT cs.agent_id) FILTER (WHERE (cs.agent_id IS NOT NULL)) AS operating_agents
   FROM ((public.companies c
     LEFT JOIN public.campaigns ca ON ((ca.company_id = c.id)))
     LEFT JOIN public.cases cs ON ((cs.company_id = c.id)))
  GROUP BY c.id, c.name
  ORDER BY c.id;


--
-- Name: vw_non_agent_users; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_non_agent_users AS
 SELECT u.id,
    u.email,
    u.full_name,
    u.phone,
    u.is_active,
    u.profile_image_url
   FROM (public.users u
     LEFT JOIN public.agents a ON ((a.user_id = u.id)))
  WHERE (a.user_id IS NULL);


--
-- Name: vw_payment_receipt_cases; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_payment_receipt_cases AS
 SELECT c.id AS case_id,
    c.sender_id,
    rr.status AS receipt_status,
    epc.harmony_state,
    epc.status AS erp_status,
    (rr.reference_number)::text AS reference_number
   FROM ((public.receipt_results rr
     JOIN public.cases c ON ((c.id = rr.case_id)))
     JOIN public.erp_payment_confirmation epc ON ((epc.reference_number = (rr.reference_number)::text)))
  WHERE ((rr.status)::text <> 'closed'::text);


--
-- Name: vw_payment_validations; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_payment_validations AS
 SELECT c.id AS case_id,
    c.sender_id,
    rr.bank_name,
    rr.transaction_type,
    rr.reference_number AS receipt_reference_number,
    rr.date AS receipt_date,
    rr."time" AS receipt_time,
    rr.amount AS receipt_amount,
    rr.amount_sent,
    rr.sender_name,
    rr.raw_text,
    rr.created_at AS ocr_date,
    erp.status AS erp_status,
    erp.harmony_state,
    erp.reference_number AS erp_reference_number,
    erp.payment_date,
    erp.payment_time,
    erp.amount AS erp_amount,
    erp.client_document,
    erp.id AS erp_id,
    erp.client_name,
    erp.contract_id
   FROM ((public.cases c
     JOIN public.receipt_results rr ON ((c.id = rr.case_id)))
     LEFT JOIN public.erp_payment_confirmation erp ON (((rr.reference_number)::text = erp.reference_number)));


--
-- Name: vw_role_permissions; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_role_permissions AS
 SELECT r.id AS role_id,
    p.id AS permission_id,
    p.code,
    p.description,
        CASE
            WHEN (rp.permission_id IS NOT NULL) THEN true
            ELSE false
        END AS assigned
   FROM ((public.roles r
     CROSS JOIN public.permissions p)
     LEFT JOIN public.role_permissions rp ON (((rp.role_id = r.id) AND (rp.permission_id = p.id))));


--
-- Name: vw_unreconciled_erp_payments; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_unreconciled_erp_payments AS
 SELECT erp.id_record,
    erp.id AS erp_id,
    erp.status AS erp_status,
    erp.harmony_state,
    erp.bank_name,
    erp.reference_number,
    erp.payment_date,
    erp.payment_time,
    erp.amount AS erp_amount,
    erp.client_document,
    erp.client_name,
    erp.contract_id,
    erp.receipt_type,
    erp.created_at
   FROM (public.erp_payment_confirmation erp
     LEFT JOIN public.receipt_results rr ON (((erp.reference_number = (rr.reference_number)::text) OR (rr.raw_text ~~ (('%'::text || erp.reference_number) || '%'::text)))))
  WHERE (rr.id IS NULL);


--
-- Name: vw_user_effective_permissions; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_user_effective_permissions AS
 SELECT u.id AS user_id,
    u.email AS user_email,
    uc.company_id,
    r.id AS role_id,
    r.name AS role_name,
    p.id AS permission_id,
    p.code AS permission_code,
    p.description AS permission_description
   FROM ((((public.users u
     JOIN public.user_company_roles uc ON ((uc.user_id = u.id)))
     JOIN public.roles r ON ((r.id = uc.role_id)))
     JOIN public.role_permissions rp ON ((rp.role_id = r.id)))
     JOIN public.permissions p ON ((p.id = rp.permission_id)));


--
-- Name: vw_user_permissions; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_user_permissions AS
 SELECT ucr.user_id,
    ucr.company_id,
    p.id AS permission_id,
    p.code AS permission_code,
    p.description AS permission_description,
    bool_or((rp.permission_id IS NOT NULL)) AS has_permission
   FROM (((public.user_company_roles ucr
     JOIN public.roles r ON ((r.id = ucr.role_id)))
     CROSS JOIN public.permissions p)
     LEFT JOIN public.role_permissions rp ON (((rp.role_id = r.id) AND (rp.permission_id = p.id))))
  GROUP BY ucr.user_id, ucr.company_id, p.id, p.code, p.description;


--
-- Name: vw_user_role_permissions; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_user_role_permissions AS
 SELECT ucr.company_id,
    u.id AS user_id,
    COALESCE(u.full_name, u.email) AS user_name,
    r.id AS role_id,
    r.name AS role_name,
    r.is_agent,
    p.id AS permission_id,
    p.code AS permission_code,
    p.description AS permission_description,
        CASE
            WHEN (rp.permission_id IS NOT NULL) THEN true
            ELSE false
        END AS has_permission
   FROM ((((public.user_company_roles ucr
     JOIN public.users u ON ((u.id = ucr.user_id)))
     JOIN public.roles r ON ((r.id = ucr.role_id)))
     CROSS JOIN public.permissions p)
     LEFT JOIN public.role_permissions rp ON (((rp.role_id = r.id) AND (rp.permission_id = p.id))));


--
-- Name: vw_users_companies; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.vw_users_companies AS
 SELECT DISTINCT u.id AS user_id,
    u.email,
    u.full_name,
    u.phone,
    u.is_active,
    c.id AS company_id,
    c.name AS company_name,
    u.is_super_user
   FROM ((public.users u
     JOIN public.user_company_roles ucr ON ((u.id = ucr.user_id)))
     JOIN public.companies c ON ((ucr.company_id = c.id)));


--
-- Name: whatsapp_message_control; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.whatsapp_message_control (
    id bigint NOT NULL,
    ws_message__id text NOT NULL
);


--
-- Name: whatsapp_message_control_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.whatsapp_message_control_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: whatsapp_message_control_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.whatsapp_message_control_id_seq OWNED BY public.whatsapp_message_control.id;


--
-- Name: agent_department_assignments id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.agent_department_assignments ALTER COLUMN id SET DEFAULT nextval('public.agent_department_assignments_id_seq'::regclass);


--
-- Name: campaign_whatsapp_push id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.campaign_whatsapp_push ALTER COLUMN id SET DEFAULT nextval('public.campaign_whatsapp_push_id_seq'::regclass);


--
-- Name: campaign_whatsapp_push_leads id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.campaign_whatsapp_push_leads ALTER COLUMN id SET DEFAULT nextval('public.campaign_whatsapp_push_leads_id_seq'::regclass);


--
-- Name: campaigns id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.campaigns ALTER COLUMN id SET DEFAULT nextval('public.campaigns_id_seq'::regclass);


--
-- Name: cantons id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cantons ALTER COLUMN id SET DEFAULT nextval('public.cantons_id_seq'::regclass);


--
-- Name: case_funnel id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.case_funnel ALTER COLUMN id SET DEFAULT nextval('public.case_funnel_id_seq'::regclass);


--
-- Name: case_items id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.case_items ALTER COLUMN id SET DEFAULT nextval('public.case_items_id_seq'::regclass);


--
-- Name: case_notes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.case_notes ALTER COLUMN id SET DEFAULT nextval('public.case_notes_id_seq'::regclass);


--
-- Name: cases id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cases ALTER COLUMN id SET DEFAULT nextval('public.cases_id_seq'::regclass);


--
-- Name: channel_agent_client id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channel_agent_client ALTER COLUMN id SET DEFAULT nextval('public.channel_agent_client_id_seq'::regclass);


--
-- Name: channel_integrations id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channel_integrations ALTER COLUMN id SET DEFAULT nextval('public.channel_integrations_id_seq'::regclass);


--
-- Name: channel_whatsapp_template id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channel_whatsapp_template ALTER COLUMN id SET DEFAULT nextval('public.channel_whatsapp_template_id_seq'::regclass);


--
-- Name: channels id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channels ALTER COLUMN id SET DEFAULT nextval('public.channels_id_seq'::regclass);


--
-- Name: client_campaigns id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_campaigns ALTER COLUMN id SET DEFAULT nextval('public.client_campaigns_id_seq'::regclass);


--
-- Name: client_files id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_files ALTER COLUMN id SET DEFAULT nextval('public.client_files_id_seq'::regclass);


--
-- Name: client_notes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_notes ALTER COLUMN id SET DEFAULT nextval('public.client_notes_id_seq'::regclass);


--
-- Name: client_social_accounts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_social_accounts ALTER COLUMN id SET DEFAULT nextval('public.client_social_accounts_id_seq'::regclass);


--
-- Name: clients id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.clients ALTER COLUMN id SET DEFAULT nextval('public.clients_id_seq'::regclass);


--
-- Name: companies id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.companies ALTER COLUMN id SET DEFAULT nextval('public.companies_id_seq'::regclass);


--
-- Name: countries id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.countries ALTER COLUMN id SET DEFAULT nextval('public.countries_id_seq'::regclass);


--
-- Name: custom_field_definitions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.custom_field_definitions ALTER COLUMN id SET DEFAULT nextval('public.custom_field_definitions_id_seq'::regclass);


--
-- Name: custom_field_values id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.custom_field_values ALTER COLUMN id SET DEFAULT nextval('public.custom_field_values_id_seq'::regclass);


--
-- Name: custom_list_definitions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.custom_list_definitions ALTER COLUMN id SET DEFAULT nextval('public.custom_list_definitions_id_seq'::regclass);


--
-- Name: custom_list_entity_value id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.custom_list_entity_value ALTER COLUMN id SET DEFAULT nextval('public.custom_list_entity_value_id_seq'::regclass);


--
-- Name: custom_list_values id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.custom_list_values ALTER COLUMN id SET DEFAULT nextval('public.custom_list_values_id_seq'::regclass);


--
-- Name: departments id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.departments ALTER COLUMN id SET DEFAULT nextval('public.departments_id_seq'::regclass);


--
-- Name: districts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.districts ALTER COLUMN id SET DEFAULT nextval('public.districts_id_seq'::regclass);


--
-- Name: erp_payment_confirmation id_record; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.erp_payment_confirmation ALTER COLUMN id_record SET DEFAULT nextval('public.erp_payment_confirmation_id_record_seq'::regclass);


--
-- Name: funnel_stages id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.funnel_stages ALTER COLUMN id SET DEFAULT nextval('public.funnel_stages_id_seq'::regclass);


--
-- Name: funnels id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.funnels ALTER COLUMN id SET DEFAULT nextval('public.funnels_id_seq'::regclass);


--
-- Name: integration_templates id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.integration_templates ALTER COLUMN id SET DEFAULT nextval('public.integration_templates_id_seq'::regclass);


--
-- Name: items id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.items ALTER COLUMN id SET DEFAULT nextval('public.items_id_seq'::regclass);


--
-- Name: message_status id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.message_status ALTER COLUMN id SET DEFAULT nextval('public.message_status_id_seq'::regclass);


--
-- Name: message_templates id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.message_templates ALTER COLUMN id SET DEFAULT nextval('public.message_templates_id_seq'::regclass);


--
-- Name: messages id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.messages ALTER COLUMN id SET DEFAULT nextval('public.messages_id_seq'::regclass);


--
-- Name: permissions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permissions ALTER COLUMN id SET DEFAULT nextval('public.permissions_id_seq'::regclass);


--
-- Name: provinces id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.provinces ALTER COLUMN id SET DEFAULT nextval('public.provinces_id_seq'::regclass);


--
-- Name: qr_leads id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.qr_leads ALTER COLUMN id SET DEFAULT nextval('public.qr_leads_id_seq'::regclass);


--
-- Name: receipt_process id_record; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.receipt_process ALTER COLUMN id_record SET DEFAULT nextval('public.receipt_process_id_record_seq'::regclass);


--
-- Name: receipt_results id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.receipt_results ALTER COLUMN id SET DEFAULT nextval('public.receipt_results_id_seq'::regclass);


--
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- Name: seats_sale id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.seats_sale ALTER COLUMN id SET DEFAULT nextval('public.seats_sale_id_seq'::regclass);


--
-- Name: settings id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.settings ALTER COLUMN id SET DEFAULT nextval('public.settings_id_seq'::regclass);


--
-- Name: user_channel_permissions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_channel_permissions ALTER COLUMN id SET DEFAULT nextval('public.user_channel_permissions_id_seq'::regclass);


--
-- Name: user_company_roles id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_company_roles ALTER COLUMN id SET DEFAULT nextval('public.user_company_roles_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: whatsapp_message_control id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.whatsapp_message_control ALTER COLUMN id SET DEFAULT nextval('public.whatsapp_message_control_id_seq'::regclass);


--
-- Name: agent_department_assignments agent_department_assignments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.agent_department_assignments
    ADD CONSTRAINT agent_department_assignments_pkey PRIMARY KEY (id);


--
-- Name: agents agents_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.agents
    ADD CONSTRAINT agents_pkey PRIMARY KEY (user_id);


--
-- Name: campaign_whatsapp_push_leads campaign_whatsapp_push_leads_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.campaign_whatsapp_push_leads
    ADD CONSTRAINT campaign_whatsapp_push_leads_pkey PRIMARY KEY (id);


--
-- Name: campaign_whatsapp_push campaign_whatsapp_push_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.campaign_whatsapp_push
    ADD CONSTRAINT campaign_whatsapp_push_pkey PRIMARY KEY (id);


--
-- Name: campaigns campaigns_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.campaigns
    ADD CONSTRAINT campaigns_pkey PRIMARY KEY (id);


--
-- Name: cantons cantons_country_code_code_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cantons
    ADD CONSTRAINT cantons_country_code_code_key UNIQUE (country_code, code);


--
-- Name: cantons cantons_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cantons
    ADD CONSTRAINT cantons_pkey PRIMARY KEY (id);


--
-- Name: case_funnel case_funnel_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.case_funnel
    ADD CONSTRAINT case_funnel_pkey PRIMARY KEY (id);


--
-- Name: case_items case_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.case_items
    ADD CONSTRAINT case_items_pkey PRIMARY KEY (id);


--
-- Name: case_notes case_notes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.case_notes
    ADD CONSTRAINT case_notes_pkey PRIMARY KEY (id);


--
-- Name: cases cases_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cases
    ADD CONSTRAINT cases_pkey PRIMARY KEY (id);


--
-- Name: channel_agent_client channel_agent_client_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channel_agent_client
    ADD CONSTRAINT channel_agent_client_pkey PRIMARY KEY (id);


--
-- Name: channel_integrations channel_integrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channel_integrations
    ADD CONSTRAINT channel_integrations_pkey PRIMARY KEY (id);


--
-- Name: channel_whatsapp_template channel_whatsapp_template_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channel_whatsapp_template
    ADD CONSTRAINT channel_whatsapp_template_pkey PRIMARY KEY (id);


--
-- Name: channels channels_code_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channels
    ADD CONSTRAINT channels_code_key UNIQUE (code);


--
-- Name: channels channels_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channels
    ADD CONSTRAINT channels_pkey PRIMARY KEY (id);


--
-- Name: client_campaigns client_campaigns_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_campaigns
    ADD CONSTRAINT client_campaigns_pkey PRIMARY KEY (id);


--
-- Name: client_files client_files_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_files
    ADD CONSTRAINT client_files_pkey PRIMARY KEY (id);


--
-- Name: client_notes client_notes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_notes
    ADD CONSTRAINT client_notes_pkey PRIMARY KEY (id);


--
-- Name: client_social_accounts client_social_accounts_channel_id_external_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_social_accounts
    ADD CONSTRAINT client_social_accounts_channel_id_external_id_key UNIQUE (channel_id, external_id);


--
-- Name: client_social_accounts client_social_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_social_accounts
    ADD CONSTRAINT client_social_accounts_pkey PRIMARY KEY (id);


--
-- Name: clients clients_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.clients
    ADD CONSTRAINT clients_pkey PRIMARY KEY (id);


--
-- Name: companies companies_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT companies_pkey PRIMARY KEY (id);


--
-- Name: countries countries_code_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.countries
    ADD CONSTRAINT countries_code_key UNIQUE (iso_code);


--
-- Name: countries countries_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.countries
    ADD CONSTRAINT countries_pkey PRIMARY KEY (id);


--
-- Name: custom_field_definitions custom_field_definitions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.custom_field_definitions
    ADD CONSTRAINT custom_field_definitions_pkey PRIMARY KEY (id);


--
-- Name: custom_field_values custom_field_values_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.custom_field_values
    ADD CONSTRAINT custom_field_values_pkey PRIMARY KEY (id);


--
-- Name: custom_list_definitions custom_list_definitions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.custom_list_definitions
    ADD CONSTRAINT custom_list_definitions_pkey PRIMARY KEY (id);


--
-- Name: custom_list_entity_value custom_list_entity_value_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.custom_list_entity_value
    ADD CONSTRAINT custom_list_entity_value_pkey PRIMARY KEY (id);


--
-- Name: custom_list_values custom_list_values_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.custom_list_values
    ADD CONSTRAINT custom_list_values_pkey PRIMARY KEY (id);


--
-- Name: departments departments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.departments
    ADD CONSTRAINT departments_pkey PRIMARY KEY (id);


--
-- Name: districts districts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.districts
    ADD CONSTRAINT districts_pkey PRIMARY KEY (id);


--
-- Name: erp_payment_confirmation erp_payment_confirmation_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.erp_payment_confirmation
    ADD CONSTRAINT erp_payment_confirmation_pkey PRIMARY KEY (id_record);


--
-- Name: funnel_stages funnel_stages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.funnel_stages
    ADD CONSTRAINT funnel_stages_pkey PRIMARY KEY (id);


--
-- Name: funnels funnels_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.funnels
    ADD CONSTRAINT funnels_pkey PRIMARY KEY (id);


--
-- Name: integration_templates integration_templates_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.integration_templates
    ADD CONSTRAINT integration_templates_pkey PRIMARY KEY (id);


--
-- Name: item_departments item_departments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.item_departments
    ADD CONSTRAINT item_departments_pkey PRIMARY KEY (item_id, department_id);


--
-- Name: items items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_pkey PRIMARY KEY (id);


--
-- Name: message_status message_status_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.message_status
    ADD CONSTRAINT message_status_pkey PRIMARY KEY (id);


--
-- Name: message_templates message_templates_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.message_templates
    ADD CONSTRAINT message_templates_pkey PRIMARY KEY (id);


--
-- Name: messages messages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_pkey PRIMARY KEY (id);


--
-- Name: permissions permissions_code_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_code_key UNIQUE (code);


--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- Name: provinces provinces_country_code_code_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.provinces
    ADD CONSTRAINT provinces_country_code_code_key UNIQUE (country_code, code);


--
-- Name: provinces provinces_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.provinces
    ADD CONSTRAINT provinces_pkey PRIMARY KEY (id);


--
-- Name: qr_leads qr_leads_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.qr_leads
    ADD CONSTRAINT qr_leads_pkey PRIMARY KEY (id);


--
-- Name: receipt_process receipt_process_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.receipt_process
    ADD CONSTRAINT receipt_process_pkey PRIMARY KEY (id_record);


--
-- Name: receipt_results receipt_results_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.receipt_results
    ADD CONSTRAINT receipt_results_pkey PRIMARY KEY (id);


--
-- Name: role_permissions role_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (role_id, permission_id);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: seats_sale seats_sale_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.seats_sale
    ADD CONSTRAINT seats_sale_pkey PRIMARY KEY (id);


--
-- Name: seats_sale seats_sale_seat_number_zone_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.seats_sale
    ADD CONSTRAINT seats_sale_seat_number_zone_key UNIQUE (seat_number, zone);


--
-- Name: settings settings_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.settings
    ADD CONSTRAINT settings_pkey PRIMARY KEY (id);


--
-- Name: funnel_stages uq_funnel_stage_position; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.funnel_stages
    ADD CONSTRAINT uq_funnel_stage_position UNIQUE (funnel_id, "position");


--
-- Name: integration_templates uq_it_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.integration_templates
    ADD CONSTRAINT uq_it_unique UNIQUE (integration_id, template_id);


--
-- Name: message_templates uq_mt_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.message_templates
    ADD CONSTRAINT uq_mt_unique UNIQUE (channel_id, template_name, language_code);


--
-- Name: user_channel_permissions user_channel_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_channel_permissions
    ADD CONSTRAINT user_channel_permissions_pkey PRIMARY KEY (id);


--
-- Name: user_company_roles user_company_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_company_roles
    ADD CONSTRAINT user_company_roles_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: whatsapp_message_control whatsapp_message_control_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.whatsapp_message_control
    ADD CONSTRAINT whatsapp_message_control_pkey PRIMARY KEY (id);


--
-- Name: idx_cases_company_channel; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_cases_company_channel ON public.cases USING btree (company_id, channel_id);


--
-- Name: idx_cases_company_department_status_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_cases_company_department_status_created ON public.cases USING btree (company_id, department_id, status, created_at DESC);


--
-- Name: idx_cf_case_changedat; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_cf_case_changedat ON public.case_funnel USING btree (case_id, changed_at DESC);


--
-- Name: idx_cf_case_latest; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_cf_case_latest ON public.case_funnel USING btree (case_id, to_stage_id);


--
-- Name: idx_cf_changed_by; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_cf_changed_by ON public.case_funnel USING btree (changed_by);


--
-- Name: idx_cf_funnel_stage; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_cf_funnel_stage ON public.case_funnel USING btree (funnel_id, to_stage_id);


--
-- Name: idx_channel_integrations_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_channel_integrations_active ON public.channel_integrations USING btree (company_id, channel_id) WHERE (is_active = true);


--
-- Name: idx_custom_field_values_field_entity; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_custom_field_values_field_entity ON public.custom_field_values USING btree (field_id, entity_id);


--
-- Name: idx_epc_created_at_reference; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_epc_created_at_reference ON public.erp_payment_confirmation USING btree (created_at, reference_number);


--
-- Name: idx_epc_harmony_state_with_receipt; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_epc_harmony_state_with_receipt ON public.erp_payment_confirmation USING btree (harmony_state) WHERE ((receipt_base64 IS NOT NULL) AND (receipt_base64 <> ''::text));


--
-- Name: idx_epc_state_empty_receipt; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_epc_state_empty_receipt ON public.erp_payment_confirmation USING btree (harmony_state, id) WHERE ((receipt_base64 IS NULL) OR (receipt_base64 = ''::text));


--
-- Name: idx_erp_payment_confirmation_reference; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_erp_payment_confirmation_reference ON public.erp_payment_confirmation USING btree (reference_number);


--
-- Name: idx_funnel_stages_funnel; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_funnel_stages_funnel ON public.funnel_stages USING btree (funnel_id);


--
-- Name: idx_funnel_stages_position; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_funnel_stages_position ON public.funnel_stages USING btree (funnel_id, "position");


--
-- Name: idx_messages_case_client; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_messages_case_client ON public.messages USING btree (case_id) WHERE (sender_type = 'client'::text);


--
-- Name: idx_messages_case_id_id_desc; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_messages_case_id_id_desc ON public.messages USING btree (case_id, id DESC);


--
-- Name: idx_messages_case_unread; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_messages_case_unread ON public.messages USING btree (case_id) WHERE (message_read = false);


--
-- Name: idx_mv_cases_company_department_status_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_mv_cases_company_department_status_created ON public.mv_cases_with_channels USING btree (company_id, department_id, status, created_at DESC);


--
-- Name: idx_mv_cases_unique; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_mv_cases_unique ON public.mv_cases_with_channels USING btree (case_id);


--
-- Name: idx_permissions_code; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_permissions_code ON public.permissions USING btree (code);


--
-- Name: idx_qr_leads_campaign_phone; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_qr_leads_campaign_phone ON public.qr_leads USING btree (campaign_id, contact_phone);


--
-- Name: idx_qr_leads_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_qr_leads_status ON public.qr_leads USING btree (status);


--
-- Name: idx_receipt_results_case_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_receipt_results_case_id ON public.receipt_results USING btree (case_id);


--
-- Name: idx_role_permissions_role; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_role_permissions_role ON public.role_permissions USING btree (role_id, permission_id);


--
-- Name: idx_ucr_user_company; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ucr_user_company ON public.user_company_roles USING btree (user_id, company_id);


--
-- Name: funnel_stages trg_funnel_stages_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_funnel_stages_updated_at BEFORE UPDATE ON public.funnel_stages FOR EACH ROW EXECUTE FUNCTION public.set_updated_at_column();


--
-- Name: funnels trg_funnels_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_funnels_updated_at BEFORE UPDATE ON public.funnels FOR EACH ROW EXECUTE FUNCTION public.set_updated_at_column();


--
-- Name: receipt_results trg_receipt_results_after_insert; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_receipt_results_after_insert AFTER INSERT ON public.receipt_results FOR EACH ROW EXECUTE FUNCTION public.fn_update_case_payment_found();


--
-- Name: receipt_results trg_update_erp_payment_confirmation; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_update_erp_payment_confirmation AFTER INSERT ON public.receipt_results FOR EACH ROW EXECUTE FUNCTION public.fn_update_erp_payment_confirmation();


--
-- Name: agent_department_assignments agent_department_assignments_agent_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.agent_department_assignments
    ADD CONSTRAINT agent_department_assignments_agent_id_fkey FOREIGN KEY (agent_id) REFERENCES public.agents(user_id);


--
-- Name: agent_department_assignments agent_department_assignments_department_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.agent_department_assignments
    ADD CONSTRAINT agent_department_assignments_department_id_fkey FOREIGN KEY (department_id) REFERENCES public.departments(id);


--
-- Name: agents agents_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.agents
    ADD CONSTRAINT agents_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: campaigns campaigns_company_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.campaigns
    ADD CONSTRAINT campaigns_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- Name: campaigns campaigns_funnel_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.campaigns
    ADD CONSTRAINT campaigns_funnel_id_fkey FOREIGN KEY (funnel_id) REFERENCES public.funnels(id);


--
-- Name: cantons cantons_country_code_province_code_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cantons
    ADD CONSTRAINT cantons_country_code_province_code_fkey FOREIGN KEY (country_code, province_code) REFERENCES public.provinces(country_code, code) ON DELETE CASCADE;


--
-- Name: case_items case_items_case_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.case_items
    ADD CONSTRAINT case_items_case_id_fkey FOREIGN KEY (case_id) REFERENCES public.cases(id) ON DELETE CASCADE;


--
-- Name: case_items case_items_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.case_items
    ADD CONSTRAINT case_items_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: case_items case_items_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.case_items
    ADD CONSTRAINT case_items_item_id_fkey FOREIGN KEY (item_id) REFERENCES public.items(id) ON DELETE CASCADE;


--
-- Name: case_notes case_notes_author_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.case_notes
    ADD CONSTRAINT case_notes_author_id_fkey FOREIGN KEY (author_id) REFERENCES public.users(id);


--
-- Name: case_notes case_notes_case_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.case_notes
    ADD CONSTRAINT case_notes_case_id_fkey FOREIGN KEY (case_id) REFERENCES public.cases(id);


--
-- Name: cases cases_agent_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cases
    ADD CONSTRAINT cases_agent_id_fkey FOREIGN KEY (agent_id) REFERENCES public.users(id);


--
-- Name: cases cases_campaign_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cases
    ADD CONSTRAINT cases_campaign_id_fkey FOREIGN KEY (campaign_id) REFERENCES public.campaigns(id);


--
-- Name: cases cases_client_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cases
    ADD CONSTRAINT cases_client_id_fkey FOREIGN KEY (client_id) REFERENCES public.clients(id);


--
-- Name: cases cases_company_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cases
    ADD CONSTRAINT cases_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- Name: cases cases_current_stage_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cases
    ADD CONSTRAINT cases_current_stage_id_fkey FOREIGN KEY (current_stage_id) REFERENCES public.funnel_stages(id);


--
-- Name: cases cases_department_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cases
    ADD CONSTRAINT cases_department_id_fkey FOREIGN KEY (department_id) REFERENCES public.departments(id);


--
-- Name: cases cases_funnel_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cases
    ADD CONSTRAINT cases_funnel_id_fkey FOREIGN KEY (funnel_id) REFERENCES public.funnels(id);


--
-- Name: channel_integrations channel_integrations_channel_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channel_integrations
    ADD CONSTRAINT channel_integrations_channel_id_fkey FOREIGN KEY (channel_id) REFERENCES public.channels(id);


--
-- Name: channel_integrations channel_integrations_company_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channel_integrations
    ADD CONSTRAINT channel_integrations_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- Name: client_campaigns client_campaigns_assigned_agent_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_campaigns
    ADD CONSTRAINT client_campaigns_assigned_agent_id_fkey FOREIGN KEY (assigned_agent_id) REFERENCES public.users(id);


--
-- Name: client_campaigns client_campaigns_campaign_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_campaigns
    ADD CONSTRAINT client_campaigns_campaign_id_fkey FOREIGN KEY (campaign_id) REFERENCES public.campaigns(id);


--
-- Name: client_campaigns client_campaigns_client_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_campaigns
    ADD CONSTRAINT client_campaigns_client_id_fkey FOREIGN KEY (client_id) REFERENCES public.clients(id);


--
-- Name: client_campaigns client_campaigns_funnel_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_campaigns
    ADD CONSTRAINT client_campaigns_funnel_id_fkey FOREIGN KEY (funnel_id) REFERENCES public.funnels(id);


--
-- Name: client_files client_files_client_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_files
    ADD CONSTRAINT client_files_client_id_fkey FOREIGN KEY (client_id) REFERENCES public.clients(id);


--
-- Name: client_files client_files_uploader_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_files
    ADD CONSTRAINT client_files_uploader_id_fkey FOREIGN KEY (uploader_id) REFERENCES public.users(id);


--
-- Name: client_notes client_notes_author_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_notes
    ADD CONSTRAINT client_notes_author_id_fkey FOREIGN KEY (author_id) REFERENCES public.users(id);


--
-- Name: client_notes client_notes_client_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_notes
    ADD CONSTRAINT client_notes_client_id_fkey FOREIGN KEY (client_id) REFERENCES public.clients(id);


--
-- Name: client_social_accounts client_social_accounts_channel_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_social_accounts
    ADD CONSTRAINT client_social_accounts_channel_id_fkey FOREIGN KEY (channel_id) REFERENCES public.channels(id);


--
-- Name: client_social_accounts client_social_accounts_client_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_social_accounts
    ADD CONSTRAINT client_social_accounts_client_id_fkey FOREIGN KEY (client_id) REFERENCES public.clients(id);


--
-- Name: custom_field_values custom_field_values_field_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.custom_field_values
    ADD CONSTRAINT custom_field_values_field_id_fkey FOREIGN KEY (field_id) REFERENCES public.custom_field_definitions(id) ON DELETE CASCADE;


--
-- Name: departments departments_company_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.departments
    ADD CONSTRAINT departments_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- Name: districts districts_country_code_canton_code_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.districts
    ADD CONSTRAINT districts_country_code_canton_code_fkey FOREIGN KEY (country_code, canton_code) REFERENCES public.cantons(country_code, code) ON DELETE CASCADE;


--
-- Name: channel_agent_client fk_agent; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channel_agent_client
    ADD CONSTRAINT fk_agent FOREIGN KEY (agent_id) REFERENCES public.agents(user_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: case_funnel fk_cf_case; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.case_funnel
    ADD CONSTRAINT fk_cf_case FOREIGN KEY (case_id) REFERENCES public.cases(id) ON DELETE CASCADE;


--
-- Name: case_funnel fk_cf_from; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.case_funnel
    ADD CONSTRAINT fk_cf_from FOREIGN KEY (from_stage_id) REFERENCES public.funnel_stages(id) ON DELETE SET NULL;


--
-- Name: case_funnel fk_cf_funnel; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.case_funnel
    ADD CONSTRAINT fk_cf_funnel FOREIGN KEY (funnel_id) REFERENCES public.funnels(id) ON DELETE RESTRICT;


--
-- Name: case_funnel fk_cf_to; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.case_funnel
    ADD CONSTRAINT fk_cf_to FOREIGN KEY (to_stage_id) REFERENCES public.funnel_stages(id) ON DELETE RESTRICT;


--
-- Name: case_funnel fk_cf_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.case_funnel
    ADD CONSTRAINT fk_cf_user FOREIGN KEY (changed_by) REFERENCES public.users(id) ON DELETE RESTRICT;


--
-- Name: channel_agent_client fk_channel; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channel_agent_client
    ADD CONSTRAINT fk_channel FOREIGN KEY (channel_id) REFERENCES public.channels(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: channel_agent_client fk_client; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channel_agent_client
    ADD CONSTRAINT fk_client FOREIGN KEY (client_id) REFERENCES public.clients(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: integration_templates fk_it_integration; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.integration_templates
    ADD CONSTRAINT fk_it_integration FOREIGN KEY (integration_id) REFERENCES public.channel_integrations(id);


--
-- Name: integration_templates fk_it_template; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.integration_templates
    ADD CONSTRAINT fk_it_template FOREIGN KEY (template_id) REFERENCES public.message_templates(id);


--
-- Name: message_templates fk_mt_channel; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.message_templates
    ADD CONSTRAINT fk_mt_channel FOREIGN KEY (channel_id) REFERENCES public.channels(id);


--
-- Name: funnel_stages funnel_stages_funnel_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.funnel_stages
    ADD CONSTRAINT funnel_stages_funnel_id_fkey FOREIGN KEY (funnel_id) REFERENCES public.funnels(id) ON DELETE CASCADE;


--
-- Name: item_departments item_departments_department_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.item_departments
    ADD CONSTRAINT item_departments_department_id_fkey FOREIGN KEY (department_id) REFERENCES public.departments(id);


--
-- Name: item_departments item_departments_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.item_departments
    ADD CONSTRAINT item_departments_item_id_fkey FOREIGN KEY (item_id) REFERENCES public.items(id);


--
-- Name: items items_company_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- Name: messages messages_case_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_case_id_fkey FOREIGN KEY (case_id) REFERENCES public.cases(id);


--
-- Name: provinces provinces_country_code_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.provinces
    ADD CONSTRAINT provinces_country_code_fkey FOREIGN KEY (country_code) REFERENCES public.countries(iso_code) ON DELETE CASCADE;


--
-- Name: qr_leads qr_leads_campaign_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.qr_leads
    ADD CONSTRAINT qr_leads_campaign_id_fkey FOREIGN KEY (campaign_id) REFERENCES public.campaigns(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: qr_leads qr_leads_company_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.qr_leads
    ADD CONSTRAINT qr_leads_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: qr_leads qr_leads_department_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.qr_leads
    ADD CONSTRAINT qr_leads_department_id_fkey FOREIGN KEY (department_id) REFERENCES public.departments(id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: qr_leads qr_leads_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.qr_leads
    ADD CONSTRAINT qr_leads_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: role_permissions role_permissions_permission_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES public.permissions(id);


--
-- Name: role_permissions role_permissions_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- Name: user_channel_permissions user_channel_permissions_channel_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_channel_permissions
    ADD CONSTRAINT user_channel_permissions_channel_id_fkey FOREIGN KEY (channel_id) REFERENCES public.channels(id);


--
-- Name: user_channel_permissions user_channel_permissions_company_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_channel_permissions
    ADD CONSTRAINT user_channel_permissions_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- Name: user_channel_permissions user_channel_permissions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_channel_permissions
    ADD CONSTRAINT user_channel_permissions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_company_roles user_company_roles_company_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_company_roles
    ADD CONSTRAINT user_company_roles_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- Name: user_company_roles user_company_roles_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_company_roles
    ADD CONSTRAINT user_company_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- Name: user_company_roles user_company_roles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_company_roles
    ADD CONSTRAINT user_company_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

\unrestrict rbGpHnXKUrGkQ3B0vN2jGJhpUDbrfqkscoYa29rZTivHukUvT1H24zSh6OIfMh1

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop tables here
-- +goose StatementEnd
