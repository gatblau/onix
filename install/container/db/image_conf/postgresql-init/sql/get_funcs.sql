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
DO
  $$
    BEGIN

      /*
        gets an item by its natural key.
        use: select * from item('the_item_key')
       */
      CREATE OR REPLACE FUNCTION item(key_param character varying)
        RETURNS TABLE(
          id bigint,
          key character varying,
          name character varying,
          description text,
          status smallint,
          item_type_key character varying,
          meta jsonb,
          tag text[],
          attribute hstore,
          version bigint,
          created timestamp(6) with time zone,
          updated timestamp(6) with time zone,
          changed_by character varying
          )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        RETURN QUERY
          SELECT
            i.id,
            i.key,
            i.name,
            i.description,
            i.status,
            it.key as item_type_key,
            i.meta,
            i.tag,
            i.attribute,
            i.version,
            i.created,
            i.updated,
            i.changed_by
          FROM item i
                 INNER JOIN item_type it
                            ON i.item_type_id = it.id
          WHERE i.key = key_param;
      END;
      $BODY$;

      ALTER FUNCTION item(character varying)
        OWNER TO onix;

      /*
        gets an item_type by its natural key.
        use: select * from item_type('the_item_type_key')
       */
      CREATE OR REPLACE FUNCTION item_type(key_param character varying)
        RETURNS TABLE(
          id integer,
          key character varying,
          name character varying,
          description text,
          attr_valid hstore,
          filter jsonb,
          meta_schema jsonb,
          version bigint,
          created timestamp(6) with time zone,
          updated timestamp(6) with time zone,
          changed_by character varying,
          model_key character varying
        )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        RETURN QUERY
          SELECT
            i.id,
            i.key,
            i.name,
            i.description,
            i.attr_valid,
            i.filter,
            i.meta_schema,
            i.version,
            i.created,
            i.updated,
            i.changed_by,
            m.key AS model_key
          FROM item_type i
            INNER JOIN model m
              ON i.model_id = m.id
          WHERE i.key = key_param;
      END;
      $BODY$;

      ALTER FUNCTION item_type(character varying)
        OWNER TO onix;

      /*
        gets a link by its natural key.
        use: select * from link('the_link_key')
       */
      CREATE OR REPLACE FUNCTION link(key_param character varying)
        RETURNS TABLE(
          id bigint,
          "key" character varying,
          link_type_key character varying,
          start_item_key character varying,
          end_item_key character varying,
          description text,
          meta jsonb,
          tag text[],
          attribute hstore,
          version bigint,
          created TIMESTAMP(6) WITH TIME ZONE,
          updated timestamp(6) WITH TIME ZONE,
          changed_by character varying
        )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        RETURN QUERY
          SELECT
            l.id,
            l.key,
            lt.key         as link_type_key,
            start_item.key as start_item_key,
            end_item.key   as end_item_key,
            l.description,
            l.meta,
            l.tag,
            l.attribute,
            l.version,
            l.created,
            l.updated,
            l.changed_by
          FROM link l
                 INNER JOIN link_type lt
                            ON l.link_type_id = lt.id
                 INNER JOIN item start_item
                            ON l.start_item_id = start_item.id
                 INNER JOIN item end_item
                            ON l.end_item_id = end_item.id
          WHERE l.key = key_param;
      END;
      $BODY$;

      ALTER FUNCTION link(character varying)
        OWNER TO onix;

      /*
        gets a link_type by its natural key.
        use: select * from link_type('the_link_type_key')
       */
      CREATE OR REPLACE FUNCTION link_type(key_param character varying)
        RETURNS TABLE(
          id integer,
          key character varying,
          name character varying,
          description text,
          attr_valid hstore,
          meta_schema jsonb,
          version bigint,
          created timestamp(6) with time zone,
          updated timestamp(6) with time zone,
          changed_by character varying,
          model_key character varying
        )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        RETURN QUERY
          SELECT
            l.id,
            l.key,
            l.name,
            l.description,
            l.attr_valid,
            l.meta_schema,
            l.version,
            l.created,
            l.updated,
            l.changed_by,
            m.key as model_key
          FROM link_type l
          INNER JOIN model m
            ON l.model_id = m.id
          WHERE l.key = key_param;
      END;
      $BODY$;

      ALTER FUNCTION link_type(character varying)
        OWNER TO onix;

      /*
        gets a Link_rule by its natural key.
        use: select * from link_rule('the_link_rule_key')
       */
      CREATE OR REPLACE FUNCTION link_rule(key_param character varying)
        RETURNS TABLE(
          id bigint,
          key character varying(300),
          name character varying(200),
          description text,
          link_type_key character varying,
          start_item_type_key character varying,
          end_item_type_key character varying,
          version bigint,
          created timestamp(6) with time zone,
          updated timestamp(6) with time zone,
          changed_by character varying(100)
          )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        RETURN QUERY
          SELECT
            r.id,
            r.key,
            r.name,
            r.description,
            lt.key              AS link_type_key,
            start_item_type.key AS start_item_key,
            end_item_type.key   AS end_item_type_key,
            r.version,
            r.created,
            r.updated,
            r.changed_by
          FROM link_rule r
                 INNER JOIN item_type start_item_type
                            ON r.start_item_type_id = start_item_type.id
                 INNER JOIN item_type end_item_type
                            ON r.end_item_type_id = end_item_type.id
                 INNER JOIN link_type lt
                            ON r.link_type_id = lt.id
          WHERE r.key = key_param;
      END;
      $BODY$;

      ALTER FUNCTION link_rule(character varying)
        OWNER TO onix;

      /*
          model(model_key_param): gets the model specified by the model_key_param.
          use: select * model(model_key_param)
         */
      CREATE OR REPLACE FUNCTION model(model_key_param character varying)
        RETURNS TABLE(
          id integer,
          key character varying,
          name character varying,
          description text,
          version bigint,
          created timestamp(6) with time zone,
          updated timestamp(6) with time zone,
          changed_by character varying
          )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        RETURN QUERY
          SELECT
            m.id,
            m.key,
            m.name,
            m.description,
            m.version,
            m.created,
            m.updated,
            m.changed_by
          FROM model m
          WHERE m.key = model_key_param;
      END;
      $BODY$;

      ALTER FUNCTION model(character varying)
        OWNER TO onix;

      /*
        get_models(): gets all models in the system.
        use: select * from get_models()
      */
      CREATE OR REPLACE FUNCTION get_models()
        RETURNS TABLE(
          id integer,
          key character varying,
          name character varying,
          description text,
          version bigint,
          created timestamp(6) with time zone,
          updated timestamp(6) with time zone,
          changed_by character varying
          )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        RETURN QUERY
          SELECT
            m.id,
            m.key,
            m.name,
            m.description,
            m.version,
            m.created,
            m.updated,
            m.changed_by
          FROM model m;
      END;
      $BODY$;

      ALTER FUNCTION get_models()
        OWNER TO onix;

    END
    $$;