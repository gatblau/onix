/*
    Onix CMDB - Copyright (c) 2018-2019 by www.gatblau.org

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
    Unless required by applicable law or agreed to in writing, software distributed under
    the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
    either express or implied.
    See the License for the specific language governing permissions and limitations under the License.

    Contributors to this project, hereby assign copyright in this code to the project,
    to be licensed under the same terms as the rest of the code.
*/

SELECT set_item_type(
  'SUPER_HOST',
  'An amazing host.',
  'A host type for testing purposes only.',
  'wbs=>required, bu=>allowed'::hstore, -- attribute validation to be applied to items of this type
  1, -- version
  'onix' -- changed by
)

SELECT set_item(
  'KEY01'::character varying, -- item natural key
  'A name.'::character varying,
  'A description.'::text,
  '{"key1":"value1", "key2":"value345"}'::jsonb, -- json formatted data
  array['tag1', 'tag2', 'tag3']::text[], -- tags for searching
  'wbs=>0'::hstore, -- item attributes
  0::smallint, -- a status number
  'SUPER_HOST', -- the type of item
  1, -- version
  'onix' -- changed by
)

SELECT set_item(
  'KEY02'::character varying, -- item natural key
  'Another name.'::character varying,
  'Another description.'::text,
  '{"key1":"value1", "key2":"value345"}'::jsonb, -- json formatted data
  array['tag1', 'tag2', 'tag3']::text[], -- tags for searching
  'wbs=>0, bu=>AERT'::hstore, -- item attributes
  0::smallint, -- a status number
  'SUPER_HOST', -- the type of item
  1, -- version
  'onix' -- changed by
)

SELECT set_link_type(
  'link01',
  'Test link type.',
  'A link type for testing only.',
  'wbs=>required, bu=>allowed', -- attribute validation to be applied to links of this type
  1, -- version
  'onix' -- changed by
)

SELECT set_link(
   'TEST-LINK-01', -- link natural key
   'link01', -- link type
   'KEY01', -- item 1
   'KEY02', -- item 2
   'Join two items for testing purposes.', -- description
   '{"key1":"value1", "key2":"value345"}'::jsonb, -- data in json format
   array['tag1', 'tag2']::text[], -- tags for searching
   'wbs=>1.0, bu=>2.1'::hstore,
   1, -- version
   'onix' -- changed by
)

SELECT * FROM item_type('SUPER_HOST')

SELECT * FROM item('KEY01')

SELECT * FROM link_type('link01')

SELECT * FROM link('TEST-LINK-01')

SELECT delete_link_type('link01', true)

SELECT delete_item_type('SUPER_HOST', true)
