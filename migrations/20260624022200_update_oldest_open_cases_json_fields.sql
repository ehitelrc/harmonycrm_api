-- +goose Up
-- +goose StatementBegin
DROP MATERIALIZED VIEW IF EXISTS public.vw_case_dashboard_by_company_with_department;
DROP MATERIALIZED VIEW IF EXISTS public.vw_case_dashboard_by_company;

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
             round(avg((EXTRACT(epoch FROM (cs_1.closed_at - CASE WHEN cs_1.started_at IS NULL OR cs_1.started_at <= '1970-01-02'::timestamp THEN cs_1.created_at ELSE cs_1.started_at END)) / (3600)::numeric)), 2) AS avg_close_hours
            FROM (public.cases cs_1
              LEFT JOIN public.users u ON ((u.id = cs_1.agent_id)))
           WHERE (cs_1.agent_id IS NOT NULL)
           GROUP BY cs_1.company_id, u.id, u.full_name
        ), oldest_open AS (
          SELECT cs_1.company_id,
             cs_1.id AS case_id,
             cl.full_name AS client_name,
             cl.phone AS client_phone,
             cs_1.created_at,
             (SELECT MAX(m.created_at) FROM public.messages m WHERE m.case_id = cs_1.id) AS last_message_at
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
    count(cs.id) FILTER (WHERE (
        cs.status IN ('open', 'in_progress')
        AND (
            SELECT m.sender_type
            FROM public.messages m
            WHERE m.case_id = cs.id
            ORDER BY m.id DESC
            LIMIT 1
        ) = 'client'
    )) AS unanswered_cases,
    count(cs.id) FILTER (WHERE (cs.agent_id IS NULL)) AS unassigned_agents,
    count(cs.id) FILTER (WHERE (cs.client_id IS NULL)) AS unassigned_clients,
    round(avg((EXTRACT(epoch FROM (cs.closed_at - CASE WHEN cs.started_at IS NULL OR cs.started_at <= '1970-01-02'::timestamp THEN cs.created_at ELSE cs.started_at END)) / (3600)::numeric)), 2) AS avg_close_hours,
    ( SELECT json_agg(json_build_object('channel_id', ch.channel_id, 'channel_name', ch.channel_name, 'open_cases', ch.open_cases, 'closed_cases', ch.closed_cases)) AS json_agg
           FROM channel_stats ch
          WHERE (ch.company_id = c.id)) AS cases_by_channel,
    ( SELECT json_agg(json_build_object('agent_id', a.agent_id, 'agent_name', a.agent_name, 'open_cases', a.open_cases, 'closed_cases', a.closed_cases, 'avg_close_hours', a.avg_close_hours)) AS json_agg
           FROM agent_stats a
          WHERE (a.company_id = c.id)) AS cases_by_agent,
    ( SELECT json_agg(json_build_object('case_id', o.case_id, 'client_name', o.client_name, 'client_phone', o.client_phone, 'created_at', o.created_at, 'last_message_at', o.last_message_at) ORDER BY o.created_at) AS json_agg
           FROM oldest_open o
          WHERE (o.company_id = c.id)
         LIMIT 20) AS oldest_open_cases
   FROM (public.companies c
     LEFT JOIN public.cases cs ON ((cs.company_id = c.id)))
  GROUP BY c.id, c.name
  ORDER BY c.id
  WITH DATA;

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
             round(avg((EXTRACT(epoch FROM (cs_1.closed_at - CASE WHEN cs_1.started_at IS NULL OR cs_1.started_at <= '1970-01-02'::timestamp THEN cs_1.created_at ELSE cs_1.started_at END)) / 3600.0)), 2) AS avg_close_hours
            FROM (public.cases cs_1
              LEFT JOIN public.users u ON ((u.id = cs_1.agent_id)))
           WHERE (cs_1.agent_id IS NOT NULL)
           GROUP BY cs_1.company_id, cs_1.department_id, u.id, u.full_name
        ), oldest_open AS (
          SELECT cs_1.company_id,
             cs_1.department_id,
             cs_1.id AS case_id,
             cl.full_name AS client_name,
             cl.phone AS client_phone,
             cs_1.created_at,
             (SELECT MAX(m.created_at) FROM public.messages m WHERE m.case_id = cs_1.id) AS last_message_at
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
    count(cs.id) FILTER (WHERE (
        cs.status IN ('open', 'in_progress')
        AND (
            SELECT m.sender_type
            FROM public.messages m
            WHERE m.case_id = cs.id
            ORDER BY m.id DESC
            LIMIT 1
        ) = 'client'
    )) AS unanswered_cases,
    count(cs.id) FILTER (WHERE (cs.agent_id IS NULL)) AS unassigned_agents,
    count(cs.id) FILTER (WHERE (cs.client_id IS NULL)) AS unassigned_clients,
    round(avg((EXTRACT(epoch FROM (cs.closed_at - CASE WHEN cs.started_at IS NULL OR cs.started_at <= '1970-01-02'::timestamp THEN cs.created_at ELSE cs.started_at END)) / 3600.0)), 2) AS avg_close_hours,
    ( SELECT json_agg(json_build_object('channel_id', ch.channel_id, 'channel_name', ch.channel_name, 'department_id', ch.department_id, 'open_cases', ch.open_cases, 'closed_cases', ch.closed_cases)) AS json_agg
           FROM channel_stats ch
          WHERE ((ch.company_id = c.id) AND (ch.department_id = cs.department_id))) AS cases_by_channel,
    ( SELECT json_agg(json_build_object('agent_id', a.agent_id, 'agent_name', a.agent_name, 'department_id', a.department_id, 'open_cases', a.open_cases, 'closed_cases', a.closed_cases, 'avg_close_hours', a.avg_close_hours)) AS json_agg
           FROM agent_stats a
          WHERE ((a.company_id = c.id) AND (a.department_id = cs.department_id))) AS cases_by_agent,
    ( SELECT json_agg(json_build_object('case_id', o.case_id, 'client_name', o.client_name, 'client_phone', o.client_phone, 'department_id', o.department_id, 'created_at', o.created_at, 'last_message_at', o.last_message_at) ORDER BY o.created_at) AS json_agg
           FROM oldest_open o
          WHERE ((o.company_id = c.id) AND (o.department_id = cs.department_id))
         LIMIT 20) AS oldest_open_cases
   FROM (public.companies c
     LEFT JOIN public.cases cs ON ((cs.company_id = c.id)))
  GROUP BY c.id, c.name, cs.department_id
  ORDER BY c.id
  WITH DATA;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert is handled by previous version definitions, but we can do simple drops
DROP MATERIALIZED VIEW IF EXISTS public.vw_case_dashboard_by_company_with_department;
DROP MATERIALIZED VIEW IF EXISTS public.vw_case_dashboard_by_company;
-- +goose StatementEnd
