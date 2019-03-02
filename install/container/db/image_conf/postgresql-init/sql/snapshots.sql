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
DO $$
BEGIN

/*
  create_snapshot(...):
    creates a new snapshot for a subtree starting at the specified root item.
  */
CREATE OR REPLACE FUNCTION create_snapshot(
  root_item_key_param character varying,
  snapshot_label_param character varying,
  snapshot_name_param character varying,
  snapshot_description_param text,
  snapshot_created_by_param character varying
  )
  RETURNS VOID
  LANGUAGE 'plpgsql'
  VOLATILE
AS $BODY$
DECLARE
  tree_ids BIGINT[];
  loop_item RECORD;
  item_store HSTORE = ''::HSTORE;
  link_store HSTORE = ''::HSTORE;
  root_item_id BIGINT;
BEGIN
  root_item_id := (SELECT i.id FROM item i WHERE i.key = root_item_key_param);

  tree_ids := (SELECT get_child_items(root_item_id));
  FOR loop_item IN
    SELECT DISTINCT ON (ic.id) ic.id, ic.version FROM item_change ic
    WHERE id = ANY(tree_ids::BIGINT[])
    ORDER BY ic.id ASC, ic.changed DESC
  LOOP
     item_store := item_store || hstore(loop_item.id::TEXT, loop_item.version::TEXT);
  END LOOP;

  tree_ids := (SELECT get_child_links(root_item_id));
  FOR loop_item IN
    SELECT DISTINCT ON (lc.id) lc.id, lc.version FROM link_change lc
    WHERE id = ANY(tree_ids::BIGINT[])
    ORDER BY lc.id ASC, lc.changed DESC
  LOOP
     link_store := link_store || hstore(loop_item.id::TEXT, loop_item.version::TEXT);
  END LOOP;

  INSERT INTO snapshot (
    label,
    root_item_key,
    name,
    description,
    item_data,
    link_data,
    version,
    changed_by
  )
  VALUES (
    snapshot_label_param,
    root_item_key_param,
    snapshot_name_param,
    snapshot_description_param,
    item_store,
    link_store,
    1,
    snapshot_created_by_param
  );

END;
$BODY$;

/*
  update_snapshot(...):
    updates some of the attributes of an existing snapshot (e.g. label, name and description).
  */
CREATE OR REPLACE FUNCTION update_snapshot(
  root_item_key_param character varying,
  current_label_param character varying,
  new_label_param character varying,
  name_param character varying,
  description_param text,
  changed_by_param character varying,
  local_version_param bigint
)
  RETURNS TABLE(result char(1))
  LANGUAGE 'plpgsql'
  VOLATILE
AS $BODY$
DECLARE
  result char(1); -- the result status for the upsert
  current_version bigint; -- the version of the row before the update or null if no row
  rows_affected integer;
BEGIN
  -- gets the current version
  SELECT version
  FROM snapshot
  WHERE root_item_key = root_item_key_param
    AND label = current_label_param
  INTO current_version;

  UPDATE snapshot
  SET
    label = new_label_param,
    name = name_param,
    description = description_param,
    version = version + 1,
    updated = current_timestamp,
    changed_by = changed_by_param
  WHERE root_item_key = root_item_key_param
    AND label = current_label_param
    -- concurrency management - optimistic locking (disabled if local_version_param is null)
    AND (local_version_param = current_version OR local_version_param IS NULL)
    AND (
      label != new_label_param OR
      name != name_param OR
      description != description_param
    );
  GET DIAGNOSTICS rows_affected := ROW_COUNT;
  SELECT get_update_status(current_version, local_version_param, rows_affected > 0) INTO result;
  RETURN QUERY SELECT result;
END;
$BODY$;

  /*
    delete_item
   */
  CREATE OR REPLACE FUNCTION delete_snapshot(
    root_item_key_param character varying,
    label_param character varying
  )
    RETURNS VOID
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE
  AS $BODY$
  BEGIN
    DELETE FROM snapshot
    WHERE label = label_param
    AND root_item_key = root_item_key_param;
  END
  $BODY$;

  ALTER FUNCTION delete_snapshot(character varying, character varying)
    OWNER TO onix;

END
$$;