/*GRP-GNU-AGPL******************************************************************

File: testdata.sql

Copyright (C) 2021  Team Georepublic <info@georepublic.de>

Developer(s):
Copyright (C) 2021  Ashish Kumar <ashishkr23438@gmail.com>

-----

This file is part of pg_scheduleserv.

pg_scheduleserv is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

pg_scheduleserv is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with pg_scheduleserv.  If not, see <https://www.gnu.org/licenses/>.

******************************************************************GRP-GNU-AGPL*/

BEGIN;

COPY public.projects (id, name, data, created_at, updated_at, deleted) FROM stdin;
3909655254191459782	Sample Project	"random"	2021-10-22 23:29:31.618091	2021-10-22 23:29:31.618091	f
2593982828701335033		{"s": 1}	2021-10-24 19:52:52.303672	2021-10-24 19:52:52.303672	f
8943284028902589305		{"s": 1}	2021-10-24 19:52:52.303672	2021-10-24 19:52:52.303672	f
\.


COPY public.jobs (id, location_index, service, delivery, pickup, skills, priority, project_id, data, created_at, updated_at, deleted) FROM stdin;
6362411701075685873	32234010232342	00:02:25	{10,20}	{20,30}	{5,50,100}	11	2593982828701335033	{"key": "value"}	2021-10-24 20:31:25.968634	2021-10-24 20:31:25.968634	f
2229737119501208952	1081230000120000	00:01:01	{5,6}	{7,8}	{}	0	2593982828701335033	{"data": ["value1", 2]}	2021-10-24 21:12:24.320746	2021-10-24 21:12:24.320746	f
3324729385723589729	1081230000120000	00:00:00	{}	{}	{}	0	3909655254191459782	{"s": 1}	2021-10-24 21:12:24.320746	2021-10-24 21:12:24.320746	f
\.


COPY public.jobs_time_windows (id, tw_open, tw_close, created_at, updated_at) FROM stdin;
6362411701075685873	2020-10-10 00:00:00	2020-10-10 00:00:10	2021-10-26 21:25:41.290849	2021-10-26 21:25:41.290849
6362411701075685873	2020-10-11 00:00:00	2020-10-12 00:00:00	2021-10-26 21:25:51.58709	2021-10-26 21:25:51.58709
2229737119501208952	2020-10-10 00:10:00	2020-10-10 00:10:10	2021-10-26 21:05:14.204642	2021-10-26 21:05:14.204642
\.


COPY public.shipments (id, p_location_index, p_service, d_location_index, d_service, amount, skills, priority, project_id, data, created_at, updated_at, deleted) FROM stdin;
7794682317520784480	32234010232342	00:02:25	23345800023242	00:01:00	{10}	{5,10}	3	2593982828701335033	{"key": "value"}	2021-10-26 00:00:03.080467	2021-10-26 00:00:03.080467	f
3329730179111013588	1032234010232342	00:01:01	23345800023242	00:02:03	{}	{1}	1	2593982828701335033	{"data": 1}	2021-10-26 00:04:56.045611	2021-10-26 00:04:56.045611	f
3341766951177830852	1032234010232342	00:00:01	23345800023242	00:00:03	{1,3,5,0}	{1}	1	3909655254191459782	{"s": 1}	2021-10-26 00:05:16.67102	2021-10-26 00:05:16.67102	f
\.


COPY public.shipments_time_windows (id, kind, tw_open, tw_close, created_at, updated_at) FROM stdin;
7794682317520784480	p	2020-10-10 00:00:00	2020-10-10 00:00:00	2021-10-26 20:45:31.990635	2021-10-26 20:45:31.990635
7794682317520784480	d	2020-10-10 00:00:00	2020-10-11 00:00:00	2021-10-26 20:45:31.990635	2021-10-26 20:45:31.990635
7794682317520784480	p	2020-10-10 00:00:00	2020-10-12 00:00:00	2021-10-26 20:45:31.990635	2021-10-26 20:45:31.990635
3329730179111013588	d	2020-10-10 00:00:00	2020-10-10 00:00:00	2021-10-26 20:45:31.990635	2021-10-26 20:45:31.990635
\.


COPY public.vehicles (id, start_index, end_index, capacity, skills, tw_open, tw_close, speed_factor, project_id, data, created_at, updated_at, deleted) FROM stdin;
2550908592071787332	32234010232342	23345800023242	{1,3}	{10}	2020-01-01 00:00:00	2020-01-10 03:14:07	1.5	3909655254191459782	{"key": "value"}	2021-10-26 10:46:41.193101	2021-10-26 10:46:41.193101	f
7300272137290532980	1032234010232342	23345800023242	{1,3,5,0}	{1}	2020-01-01 10:10:00	2020-01-11 03:14:07	34.25	3909655254191459782	{"s": 1}	2021-10-26 10:47:54.437549	2021-10-26 10:47:54.437549	f
150202809001685363	1032234010232342	23345800023242	{1,3,5,0}	{1}	1970-01-01 00:00:00	2038-01-19 03:14:07	34.25	2593982828701335033	{"s": 1}	2021-10-26 10:48:19.203294	2021-10-26 10:48:19.203294	f
\.


COPY public.breaks (id, vehicle_id, service, data, created_at, updated_at, deleted) FROM stdin;
4668767710686035977	2550908592071787332	00:00:01	{"key": "value"}	2021-10-26 21:24:38.181374	2021-10-26 21:24:38.181374	f
3990300682121424906	2550908592071787332	00:05:24	{"s": 1}	2021-10-26 21:24:52.204943	2021-10-26 21:24:52.204943	f
2349284092384902582	7300272137290532980	00:05:24	{}	2021-10-26 21:24:52.204943	2021-10-26 21:24:52.204943	f
\.


COPY public.breaks_time_windows (id, tw_open, tw_close, created_at, updated_at) FROM stdin;
3990300682121424906	2020-10-10 00:00:00	2020-10-10 00:00:10	2021-10-26 21:25:41.290849	2021-10-26 21:25:41.290849
3990300682121424906	2020-10-11 00:00:00	2020-10-12 00:00:00	2021-10-26 21:25:51.58709	2021-10-26 21:25:51.58709
\.

END;
