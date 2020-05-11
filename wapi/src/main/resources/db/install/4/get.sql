/*
    Onix Config Manager - Copyright (c) 2018-2019 by www.gatblau.org

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
      CREATE OR REPLACE FUNCTION ox_item(key_param character varying,
                                         role_key_param character varying[])
        RETURNS TABLE
                (
                  id            bigint,
                  key           character varying,
                  name          character varying,
                  description   text,
                  status        smallint,
                  item_type_key character varying,
                  meta          jsonb,
                  meta_enc      boolean,
                  txt           text,
                  txt_enc       boolean,
                  enc_key_ix    smallint,
                  tag           text[],
                  attribute     hstore,
                  version       bigint,
                  created       timestamp(6) with time zone,
                  updated       timestamp(6) with time zone,
                  changed_by    character varying,
                  partition_key character varying
                )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        RETURN QUERY
          SELECT i.id,
                 i.key,
                 i.name,
                 i.description,
                 i.status,
                 it.key as item_type_key,
                 i.meta,
                 i.meta_enc,
                 i.txt,
                 i.txt_enc,
                 i.enc_key_ix,
                 i.tag,
                 i.attribute,
                 i.version,
                 i.created,
                 i.updated,
                 i.changed_by,
                 p.key  as partition_key
          FROM item i
                 INNER JOIN item_type it ON i.item_type_id = it.id
                 INNER JOIN partition p ON i.partition_id = p.id
                 INNER JOIN privilege pr ON p.id = pr.partition_id
                 INNER JOIN role r ON pr.role_id = r.id
          WHERE i.key = key_param
            AND pr.can_read = TRUE
            AND r.key = ANY (role_key_param);
      END;
      $BODY$;

      ALTER FUNCTION ox_item(
        character varying, -- key_param
        character varying[] -- role_key_param
        )
        OWNER TO onix;

      /*
        gets an item_type by its natural key.
        use: select * from item_type('the_item_type_key')
       */
      CREATE OR REPLACE FUNCTION ox_item_type(key_param character varying, role_key_param character varying[])
        RETURNS TABLE
        (
            id            integer,
            key           character varying,
            name          character varying,
            description   text,
            filter        jsonb,
            meta_schema   jsonb,
            version       bigint,
            created       timestamp(6) with time zone,
            updated       timestamp(6) with time zone,
            changed_by    character varying,
            model_key     character varying,
            root          boolean,
            notify_change char,
            tag           text[],
            encrypt_meta  boolean,
            encrypt_txt   boolean,
            style         jsonb
        )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        RETURN QUERY
          SELECT i.id,
                 i.key,
                 i.name,
                 i.description,
                 i.filter,
                 i.meta_schema,
                 i.version,
                 i.created,
                 i.updated,
                 i.changed_by,
                 m.key AS model_key,
                 k.root,
                 i.notify_change,
                 i.tag,
                 i.encrypt_meta,
                 i.encrypt_txt,
                 i.style
          FROM item_type i
                 INNER JOIN model m ON i.model_id = m.id
                 INNER JOIN privilege pr on m.partition_id = pr.partition_id
                 INNER JOIN role r on pr.role_id = r.id
            -- works out if it is a root item below
                 LEFT OUTER JOIN (
            SELECT t.id, (t.id IS NOT NULL) AS root
            FROM (
                   SELECT it.*
                   FROM item_type it
                     EXCEPT
                   SELECT it.*
                   FROM item_type it
                          INNER JOIN link_rule r
                                     ON it.id = r.end_item_type_id
                 ) AS t
          ) AS k
                                 ON k.id = i.id
          WHERE i.key = key_param
            AND pr.can_read = TRUE
            AND r.key = ANY (role_key_param);
      END;
      $BODY$;

      ALTER FUNCTION ox_item_type(character varying, character varying[])
        OWNER TO onix;

      /*
        gets a link by its natural key.
        use: select * from link('the_link_key')
       */
      CREATE OR REPLACE FUNCTION ox_link(key_param character varying, role_key_param character varying[])
        RETURNS TABLE
        (
          id             bigint,
          "key"          character varying,
          link_type_key  character varying,
          start_item_key character varying,
          end_item_key   character varying,
          description    text,
          meta           jsonb,
          meta_enc       boolean,
          txt            text,
          txt_enc        boolean,
          enc_key_ix     smallint,
          tag            text[],
          attribute      hstore,
          version        bigint,
          created        TIMESTAMP(6) WITH TIME ZONE,
          updated        timestamp(6) WITH TIME ZONE,
          changed_by     character varying,
          encrypt_txt        boolean,
          encrypt_meta       boolean
        )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        RETURN QUERY
          SELECT l.id,
                 l.key,
                 lt.key         as link_type_key,
                 start_item.key as start_item_key,
                 end_item.key   as end_item_key,
                 l.description,
                 l.meta,
                 l.meta_enc,
                 l.txt,
                 l.txt_enc,
                 l.enc_key_ix,
                 l.tag,
                 l.attribute,
                 l.version,
                 l.created,
                 l.updated,
                 l.changed_by,
                 lt.encrypt_txt,
                 lt.encrypt_meta
          FROM link l
                 INNER JOIN link_type lt ON l.link_type_id = lt.id
                 INNER JOIN item start_item ON l.start_item_id = start_item.id
                 INNER JOIN item end_item ON l.end_item_id = end_item.id
                 INNER JOIN model m on lt.model_id = m.id
                 INNER JOIN partition p on m.partition_id = p.id
                 INNER JOIN privilege pr on p.id = pr.partition_id
                 INNER JOIN role r on pr.role_id = r.id
          WHERE l.key = key_param
            AND r.key = ANY (role_key_param)
            AND pr.can_read = TRUE;
      END;
      $BODY$;

      ALTER FUNCTION ox_link(character varying, character varying[])
        OWNER TO onix;

      /*
        gets a link_type by its natural key.
        use: select * from ox_link_type('the_link_type_key')
       */
      CREATE OR REPLACE FUNCTION ox_link_type(key_param character varying,
                                              role_key_param character varying[])
        RETURNS TABLE
        (
            id           integer,
            key          character varying,
            name         character varying,
            description  text,
            meta_schema  jsonb,
            tag          text[],
            encrypt_meta boolean,
            encrypt_txt  boolean,
            style        jsonb,
            version      bigint,
            created      timestamp(6) with time zone,
            updated      timestamp(6) with time zone,
            changed_by   character varying,
            model_key    character varying
        )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        RETURN QUERY
          SELECT lt.id,
                 lt.key,
                 lt.name,
                 lt.description,
                 lt.meta_schema,
                 lt.tag,
                 lt.encrypt_meta,
                 lt.encrypt_txt,
                 lt.style,
                 lt.version,
                 lt.created,
                 lt.updated,
                 lt.changed_by,
                 m.key as model_key
          FROM link_type lt
                 INNER JOIN model m ON lt.model_id = m.id
                 INNER JOIN privilege pr on m.partition_id = pr.partition_id
                 INNER JOIN role r on pr.role_id = r.id
          WHERE lt.key = key_param
            AND pr.can_read = TRUE
            AND r.key = ANY (role_key_param);
      END;
      $BODY$;

      ALTER FUNCTION ox_link_type(character varying, character varying[])
        OWNER TO onix;

      /*
        ox_item_type_attribute: get a type attribute based on the specified key.
      */
      CREATE OR REPLACE FUNCTION ox_item_type_attribute(
          item_type_key_param character varying,
          type_key_param character varying,
          role_key_param character varying[]
      )
          RETURNS TABLE
          (
               id            integer,
               key           character varying,
               name          character varying,
               description   text,
               item_type_key character varying,
               type          character varying,
               def_value     character varying,
               required      boolean,
               regex         varchar,
               version       bigint,
               created       timestamp(6) with time zone,
               updated       timestamp(6) with time zone,
               changed_by    character varying
           )
          LANGUAGE 'plpgsql'
          COST 100
          STABLE
      AS $BODY$
      BEGIN
          RETURN QUERY
              SELECT ta.id,
                     ta.key,
                     ta.name,
                     ta.description,
                     it.key as item_type_key,
                     ta.type,
                     ta.def_value,
                     ta.required,
                     ta.regex,
                     ta.version,
                     ta.created,
                     ta.updated,
                     ta.changed_by
              FROM type_attribute ta
                   INNER JOIN item_type it ON ta.item_type_id = it.id
                   INNER JOIN model m ON it.model_id = m.id
                   INNER JOIN partition p ON m.partition_id = p.id
                   INNER JOIN privilege pr on p.id = pr.partition_id
                   INNER JOIN role r on pr.role_id = r.id
              WHERE it.key = item_type_key_param
                AND ta.key = type_key_param
                AND pr.can_read = TRUE
                AND r.key = ANY(role_key_param);
      END;
      $BODY$;

      ALTER FUNCTION ox_item_type_attribute(
          character varying,
          character varying,
          character varying[]
        )
          OWNER TO onix;

      /*
        ox_link_type_attribute: get a type attribute based on the specified key.
      */
      CREATE OR REPLACE FUNCTION ox_link_type_attribute(
          link_type_key_param character varying,
          type_key_param character varying,
          role_key_param character varying[]
      )
          RETURNS TABLE(
               id          integer,
               key         character varying,
               name        character varying,
               description text,
               link_type_key character varying,
               type        character varying,
               def_value   character varying,
               required    boolean,
               regex       varchar,
               version     bigint,
               created     timestamp(6) with time zone,
               updated     timestamp(6) with time zone,
               changed_by  character varying
           )
          LANGUAGE 'plpgsql'
          COST 100
          STABLE
      AS $BODY$
      BEGIN
          RETURN QUERY
              SELECT ta.id,
                     ta.key,
                     ta.name,
                     ta.description,
                     lt.key as link_type_key,
                     ta.type,
                     ta.def_value,
                     ta.required,
                     ta.regex,
                     ta.version,
                     ta.created,
                     ta.updated,
                     ta.changed_by
              FROM type_attribute ta
                       INNER JOIN link_type lt ON ta.link_type_id = lt.id
                       INNER JOIN model m ON lt.model_id = m.id
                       INNER JOIN partition p ON m.partition_id = p.id
                       INNER JOIN privilege pr on p.id = pr.partition_id
                       INNER JOIN role r on pr.role_id = r.id
              WHERE lt.key = link_type_key_param
                AND ta.key = type_key_param
                AND pr.can_read = TRUE
                AND r.key = ANY(role_key_param);
      END;
      $BODY$;

      ALTER FUNCTION ox_link_type_attribute(character varying, character varying, character varying[])
          OWNER TO onix;

      /*
        gets a Link_rule by its natural key.
        use: select * from ox_link_rule('the_link_rule_key')
       */
      CREATE OR REPLACE FUNCTION ox_link_rule(key_param character varying, role_key_param character varying[])
        RETURNS TABLE
                (
                  id                  bigint,
                  key                 character varying(300),
                  name                character varying(200),
                  description         text,
                  link_type_key       character varying,
                  start_item_type_key character varying,
                  end_item_type_key   character varying,
                  version             bigint,
                  created             timestamp(6) with time zone,
                  updated             timestamp(6) with time zone,
                  changed_by          character varying(100)
                )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        RETURN QUERY
          SELECT l.id,
                 l.key,
                 l.name,
                 l.description,
                 link_type.key       AS link_type_key,
                 start_item_type.key AS start_item_key,
                 end_item_type.key   AS end_item_type_key,
                 l.version,
                 l.created,
                 l.updated,
                 l.changed_by
          FROM link_rule l
                   INNER JOIN link_type link_type ON link_type.id = l.link_type_id
                   INNER JOIN item_type start_item_type ON start_item_type.id = l.start_item_type_id
                   INNER JOIN item_type end_item_type ON end_item_type.id = l.end_item_type_id
                   INNER JOIN model m ON link_type.model_id = m.id
                   INNER JOIN partition p ON m.partition_id = p.id
                   INNER JOIN privilege pr ON p.id = pr.partition_id
                   INNER JOIN role r ON pr.role_id = r.id
          WHERE l.key = key_param
          AND pr.can_read = TRUE
          AND r.key = ANY(role_key_param);
      END;
      $BODY$;

      ALTER FUNCTION ox_link_rule(character varying, character varying[])
        OWNER TO onix;

      /*
          user(user_key_param): gets the user specified by the user_key_param.
          use: select * ox_user(user_key_param)
         */
      CREATE OR REPLACE FUNCTION ox_user(user_key_param character varying, role_key_param character varying[])
          RETURNS TABLE
              (
                  id          bigint,
                  key         character varying,
                  name        character varying,
                  email       character varying,
                  pwd         character varying,
                  salt        character varying,
                  expires    timestamp(6) with time zone,
                  version     bigint,
                  created     timestamp(6) with time zone,
                  updated     timestamp(6) with time zone,
                  changed_by  character varying
              )
          LANGUAGE 'plpgsql'
          COST 100
          STABLE
      AS
      $BODY$
      BEGIN
          RETURN QUERY
              SELECT u.id,
                     u.key,
                     u.name,
                     u.email,
                     u.pwd,
                     u.salt,
                     u.expires,
                     u.version,
                     u.created,
                     u.updated,
                     u.changed_by
              FROM "user" u
              WHERE u.key = user_key_param
                AND ox_is_super_admin(role_key_param, FALSE);
      END;
      $BODY$;

      ALTER FUNCTION ox_user(character varying, character varying[])
          OWNER TO onix;

      /*
          ox_user_by_email(email_key_param): gets the user specified by the email address.
          use: select * from ox_user_by_email(email_key_param)
         */
      CREATE OR REPLACE FUNCTION ox_user_by_email(email_param character varying, role_key_param character varying[])
          RETURNS TABLE
                  (
                      id          bigint,
                      key         character varying,
                      name        character varying,
                      email       character varying,
                      pwd         character varying,
                      salt        character varying,
                      expires    timestamp(6) with time zone,
                      version     bigint,
                      created     timestamp(6) with time zone,
                      updated     timestamp(6) with time zone,
                      changed_by  character varying
                  )
          LANGUAGE 'plpgsql'
          COST 100
          STABLE
      AS
      $BODY$
      BEGIN
          RETURN QUERY
              SELECT u.id,
                     u.key,
                     u.name,
                     u.email,
                     u.pwd,
                     u.salt,
                     u.expires,
                     u.version,
                     u.created,
                     u.updated,
                     u.changed_by
              FROM "user" u
              WHERE u.email = email_param
                AND ox_is_super_admin(role_key_param, FALSE);
      END;
      $BODY$;

      ALTER FUNCTION ox_user_by_email(character varying, character varying[])
          OWNER TO onix;

      /*
          model(model_key_param): gets the model specified by the model_key_param.
          use: select * ox_model(model_key_param)
         */
      CREATE OR REPLACE FUNCTION ox_model(model_key_param character varying, role_key_param character varying[])
        RETURNS TABLE
                (
                  id          integer,
                  key         character varying,
                  name        character varying,
                  description text,
                  managed     boolean,
                  version     bigint,
                  created     timestamp(6) with time zone,
                  updated     timestamp(6) with time zone,
                  changed_by  character varying,
                  partition   character varying
                )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        RETURN QUERY
          SELECT m.id,
                 m.key,
                 m.name,
                 m.description,
                 m.managed,
                 m.version,
                 m.created,
                 m.updated,
                 m.changed_by,
                 p.key as partition
          FROM model m
                 INNER JOIN partition p on m.partition_id = p.id
                 INNER JOIN privilege pr on p.id = pr.partition_id
                 INNER JOIN role r on pr.role_id = r.id
          WHERE m.key = model_key_param
            AND r.key = ANY (role_key_param)
            AND pr.can_read = true;
      END;
      $BODY$;

      ALTER FUNCTION ox_model(character varying, character varying[])
        OWNER TO onix;

      /*
        get_models(): gets all models in the system.
        use: select * from ox_get_models()
      */
      CREATE OR REPLACE FUNCTION ox_get_models(role_key_param character varying[])
        RETURNS TABLE
                (
                  id          integer,
                  key         character varying,
                  name        character varying,
                  description text,
                  managed     boolean,
                  version     bigint,
                  created     timestamp(6) with time zone,
                  updated     timestamp(6) with time zone,
                  changed_by  character varying
                )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        RETURN QUERY
          SELECT m.id,
                 m.key,
                 m.name,
                 m.description,
                 m.managed,
                 m.version,
                 m.created,
                 m.updated,
                 m.changed_by
          FROM model m
                 INNER JOIN partition p on m.partition_id = p.id
                 INNER JOIN privilege pr on p.id = pr.partition_id
                 INNER JOIN role r on pr.role_id = r.id
          WHERE r.key = ANY (role_key_param)
            AND pr.can_read = TRUE;
      END;
      $BODY$;

      ALTER FUNCTION ox_get_models(character varying[])
        OWNER TO onix;

      /*
        get_users(): gets all users in the system.
        use: select * from ox_get_users()
      */
      CREATE OR REPLACE FUNCTION ox_get_users(role_key_param character varying[])
          RETURNS TABLE
              (
                  id          bigint,
                  key         character varying,
                  name        character varying,
                  email       character varying,
                  expires     timestamp(6) with time zone,
                  pwd         character varying,
                  salt        character varying,
                  version     bigint,
                  created     timestamp(6) with time zone,
                  updated     timestamp(6) with time zone,
                  changed_by  character varying
              )
          LANGUAGE 'plpgsql'
          COST 100
          STABLE
      AS
      $BODY$
      BEGIN
          -- only retrieve users if the logged user has a level 2 role
          RETURN QUERY
              SELECT u.id,
                     u.key,
                     u.name,
                     u.email,
                     u.expires,
                     u.pwd,
                     u.salt,
                     u.version,
                     u.created,
                     u.updated,
                     u.changed_by
              FROM "user" u
                WHERE ox_is_super_admin(role_key_param, FALSE);
      END;
      $BODY$;

      ALTER FUNCTION ox_get_models(character varying[])
          OWNER TO onix;

      /*
        partition(key_param, role_key_param): gets the partition specified by the key_param.
        use: select * partition(key_param, role_key_param)
       */
      CREATE OR REPLACE FUNCTION ox_partition(key_param character varying, role_key_param character varying[])
        RETURNS TABLE
                (
                  id          bigint,
                  key         character varying,
                  name        character varying,
                  description text,
                  owner       character varying,
                  version     bigint,
                  created     timestamp(6) with time zone,
                  updated     timestamp(6) with time zone,
                  changed_by  character varying
                )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        -- checks the role can modify this role
        PERFORM ox_can_manage_partition(role_key_param);

        RETURN QUERY
          SELECT p.id,
                 p.key,
                 p.name,
                 p.description,
                 p.owner,
                 p.version,
                 p.created,
                 p.updated,
                 p.changed_by
          FROM partition p
          WHERE p.key = key_param;
      END;
      $BODY$;

      ALTER FUNCTION ox_partition(character varying, character varying[])
        OWNER TO onix;

      /*
        get_partitions(): gets all partitions in the system.
        use: select * from ox_get_partitions(role_key_param)
      */
      CREATE OR REPLACE FUNCTION ox_get_partitions(role_key_param character varying[])
        RETURNS TABLE
                (
                  id          bigint,
                  key         character varying,
                  name        character varying,
                  description text,
                  owner       character varying,
                  version     bigint,
                  created     timestamp(6) with time zone,
                  updated     timestamp(6) with time zone,
                  changed_by  character varying
                )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        -- checks the role can modify this role
        PERFORM ox_can_manage_partition(role_key_param);

        RETURN QUERY
          SELECT p.id,
                 p.key,
                 p.name,
                 p.description,
                 p.owner,
                 p.version,
                 p.created,
                 p.updated,
                 p.changed_by
          FROM partition p;
      END;
      $BODY$;

      ALTER FUNCTION ox_get_partitions(character varying[])
        OWNER TO onix;

      /*
        role(key_param, role_key_param): gets the role specified by the key_param.
        use: select * role(key_param, role_key_param)
       */
      CREATE OR REPLACE FUNCTION ox_role(key_param character varying, role_key_param character varying[])
        RETURNS TABLE
                (
                  id          bigint,
                  key         character varying,
                  name        character varying,
                  description text,
                  level       integer,
                  owner       character varying,
                  version     bigint,
                  created     timestamp(6) with time zone,
                  updated     timestamp(6) with time zone,
                  changed_by  character varying
                )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        -- checks the role can modify this role
        PERFORM ox_can_manage_partition(role_key_param);

        RETURN QUERY
          SELECT r.id,
                 r.key,
                 r.name,
                 r.description,
                 r.level,
                 r.owner,
                 r.version,
                 r.created,
                 r.updated,
                 r.changed_by
          FROM role r
          WHERE r.key = key_param;
      END;
      $BODY$;

      ALTER FUNCTION ox_role(character varying, character varying[])
        OWNER TO onix;

      /*
        get_roles(): gets all roles in the system.
        use: select * from get_roles(role_key_param)
      */
      CREATE OR REPLACE FUNCTION ox_get_roles(role_key_param character varying[])
        RETURNS TABLE
                (
                  id          bigint,
                  key         character varying,
                  name        character varying,
                  description text,
                  level       integer,
                  owner       character varying,
                  version     bigint,
                  created     timestamp(6) with time zone,
                  updated     timestamp(6) with time zone,
                  changed_by  character varying
                )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      BEGIN
        -- checks the role can modify this role
        PERFORM ox_can_manage_partition(role_key_param);

        RETURN QUERY
          SELECT r.id,
                 r.key,
                 r.name,
                 r.description,
                 r.level,
                 r.owner,
                 r.version,
                 r.created,
                 r.updated,
                 r.changed_by
          FROM role r;
      END;
      $BODY$;

      ALTER FUNCTION ox_get_roles(character varying[])
        OWNER TO onix;

      CREATE OR REPLACE FUNCTION ox_get_privileges_by_role(role_key_param character varying,
                                                           logged_role_key_param character varying[])
        RETURNS TABLE
                (
                  key      character varying,
                  role_key      character varying,
                  partition_key character varying,
                  can_create    boolean,
                  can_read      boolean,
                  can_delete    boolean,
                  changed_by    character varying,
                  created       timestamp(6) with time zone
                )
        LANGUAGE 'plpgsql'
        COST 100
        STABLE
      AS
      $BODY$
      DECLARE
        role_level integer;
      BEGIN
        SELECT r.level
        FROM role r
        WHERE r.key = ANY (logged_role_key_param)
        ORDER BY r.level DESC
        LIMIT 1
          INTO role_level;

        RETURN QUERY
          SELECT pr.key,
                 r.key as role_key,
                 p.key as partition_key,
                 pr.can_create,
                 pr.can_read,
                 pr.can_delete,
                 pr.changed_by,
                 pr.created
          FROM privilege pr
                 INNER JOIN partition p ON p.id = pr.partition_id
                 INNER JOIN role r ON pr.role_id = r.id
          WHERE r.key = role_key_param
            AND ((r.owner = p.owner AND r.owner = ANY (logged_role_key_param) AND role_level = 1) OR (role_level = 2));
      END
      $BODY$;

      ALTER FUNCTION ox_get_privileges_by_role(character varying, character varying[])
        OWNER TO onix;

      /*
        ox_privilege(key_param, user_role_param): gets the privilege specified by the partition and role.
        use: select * ox_privilege(key_param, user_role_param)
       */
      CREATE OR REPLACE FUNCTION ox_privilege(
            key_param character varying,
            user_role_param character varying[]
        )
          RETURNS TABLE
          (
              key           character varying,
              role_key      character varying,
              partition_key character varying,
              can_create    boolean,
              can_read      boolean,
              can_delete    boolean,
              changed_by    character varying,
              created       timestamp(6) with time zone
          )
          LANGUAGE 'plpgsql'
          COST 100
          STABLE
      AS
      $BODY$
      DECLARE
          role_level integer;
      BEGIN
          SELECT r.level
          FROM role r
          WHERE r.key = ANY (user_role_param)
          ORDER BY r.level DESC
          LIMIT 1
          INTO role_level;

          RETURN QUERY
              SELECT pr.key,
                     r.key as role_key,
                     p.key as partition_key,
                     pr.can_create,
                     pr.can_read,
                     pr.can_delete,
                     pr.changed_by,
                     pr.created
              FROM privilege pr
                       INNER JOIN partition p ON p.id = pr.partition_id
                       INNER JOIN role r ON pr.role_id = r.id
              WHERE pr.key = key_param
                AND ((r.owner = p.owner AND r.owner = ANY (user_role_param) AND role_level = 1) OR (role_level = 2));
      END;
      $BODY$;

      ALTER FUNCTION ox_privilege(character varying, character varying[])
          OWNER TO onix;
    END
    $$;