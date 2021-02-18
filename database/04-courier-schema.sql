CREATE SCHEMA IF NOT EXISTS courier AUTHORIZATION postgres;

CREATE  TABLE "courier".audiences (
	id                   uuid NOT NULL ,
	name                 varchar(255)   ,
	created_at           timestamptz DEFAULT current_timestamp  ,
	updated_at           timestamptz DEFAULT current_timestamp  ,
	deleted_at           timestamptz  DEFAULT NULL ,
	CONSTRAINT pk_audience PRIMARY KEY ( id )
 );

CREATE  TABLE "courier".channels (
	id                   uuid NOT NULL ,
	name                 varchar(255)   ,
  audience_id          uuid NOT NULL,
	created_at           timestamptz DEFAULT current_timestamp  ,
	updated_at           timestamptz DEFAULT current_timestamp  ,
	deleted_at           timestamptz  DEFAULT NULL ,
	CONSTRAINT pk_channels PRIMARY KEY ( id )
 );

CREATE  TABLE "courier".topics (
	id                   uuid NOT NULL ,
	topic                varchar(255)   ,
  channel_id           uuid NOT NULL,
	created_at           timestamptz DEFAULT current_timestamp  ,
	updated_at           timestamptz DEFAULT current_timestamp  ,
	deleted_at           timestamptz  DEFAULT NULL ,
	CONSTRAINT pk_topics PRIMARY KEY ( id )
 );

ALTER TABLE "courier".channels ADD FOREIGN KEY (audience_id) REFERENCES "courier".audiences(id);
ALTER TABLE "courier".topics ADD FOREIGN KEY (channel_id) REFERENCES "courier".channels(id);

alter user postgres set search_path = "$user",extensions,public,courier;
