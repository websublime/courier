CREATE SCHEMA IF NOT EXISTS courier AUTHORIZATION postgres;

CREATE  TABLE "courier".audiences (
	id                   uuid NOT NULL ,
	name                 varchar(255)   ,
	created_at           timestamptz DEFAULT current_timestamp  ,
	updated_at           timestamptz DEFAULT current_timestamp  ,
	deleted_at           timestamptz  DEFAULT NULL ,
	CONSTRAINT pk_audience PRIMARY KEY ( id )
 );

CREATE  TABLE "courier".topics (
	id                   uuid NOT NULL ,
	topic                varchar(255)   ,
  audience_id          uuid NOT NULL,
	created_at           timestamptz DEFAULT current_timestamp  ,
	updated_at           timestamptz DEFAULT current_timestamp  ,
	deleted_at           timestamptz  DEFAULT NULL ,
	CONSTRAINT pk_topics PRIMARY KEY ( id )
 );

CREATE  TABLE "courier".messages (
	id                   uuid NOT NULL ,
	message              json,
  topic_id             uuid NOT NULL,
	created_at           timestamptz DEFAULT current_timestamp  ,
	updated_at           timestamptz DEFAULT current_timestamp  ,
	deleted_at           timestamptz  DEFAULT NULL ,
	CONSTRAINT pk_messages PRIMARY KEY ( id )
 );

CREATE TRIGGER on_audiences_handler
  AFTER INSERT OR UPDATE OR DELETE ON "courier".audiences 
  FOR EACH ROW EXECUTE PROCEDURE extensions.notify_hook();

CREATE TRIGGER on_topics_handler
  AFTER INSERT OR UPDATE OR DELETE ON "courier".topics 
  FOR EACH ROW EXECUTE PROCEDURE extensions.notify_hook();

ALTER TABLE "courier".topics ADD FOREIGN KEY (audience_id) REFERENCES "courier".audiences(id);
ALTER TABLE "courier".messages ADD FOREIGN KEY (topic_id) REFERENCES "courier".topics(id);

alter user postgres set search_path = "$user",extensions,public,courier;
