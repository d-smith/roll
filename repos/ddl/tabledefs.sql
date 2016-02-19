create or replace table rolldb.admin (
	name varchar(256) primary key
);

grant select, update, insert, delete
on rolldb.admin
to rolluser;

create or replace table rolldb.developer (
    email varchar(256) primary key,
    id varchar(256),
    firstName varchar(60),
    lastName varchar(60)
);

grant select, update, insert, delete
on rolldb.developer
to rolluser;

create or replace table rolldb.application (
    applicationName varchar(150) not null,
    clientId varchar(100) not null,
    developerEmail varchar(256) not null,
    developerId varchar(256) not null,
    loginProvider varchar(256) not null,
    redirectUri varchar(512) not null,
    jwtFlowAudience varchar(256),
    jwtFlowIssuer varchar(256),
    jwtFlowPublicKey varchar(2048),
    primary key(applicationName, developerEmail)
);

grant select, update, insert, delete
on rolldb.application
to rolluser;

/* TODO - add proper constraints once initial mariadb support is in place. */
