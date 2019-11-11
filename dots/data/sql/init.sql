create schema citychain collate utf8mb4_0900_ai_ci;

create table block
(
	id varchar(36) not null
		primary key,
	version bigint default 1 not null,
	created_at bigint not null,
	updated_at bigint not null,
	deleted_at timestamp null,
	done_block varchar(80) default '-1' not null
);

create index idx_block_deleted_at
	on block (deleted_at);

create table purchase
(
	id varchar(36) not null
		primary key,
	version bigint default 1 not null,
	created_at bigint not null,
	updated_at bigint not null,
	deleted_at timestamp null,
	purchase_id varchar(36) not null,
	tx_id varchar(66) not null,
	from_addr varchar(42) not null,
	to_addr varchar(42) not null,
	contract varchar(42) not null,
	nonce bigint unsigned not null,
	status varchar(16) default 'Generate' not null,
	tx_json varchar(10240) not null,
	done tinyint(1) default 0 not null,
	constraint uq_nonce
		unique (from_addr, nonce),
	constraint uq_purchase
		unique (purchase_id)
);

create index idx_purchase_deleted_at
	on purchase (deleted_at);

create table tx20
(
	id varchar(36) not null
		primary key,
	version bigint default 1 not null,
	created_at bigint not null,
	updated_at bigint not null,
	deleted_at timestamp null,
	purchase_id varchar(36) not null,
	tx_id varchar(66) not null,
	from_addr varchar(42) not null,
	to_addr varchar(42) not null,
	contract varchar(42) not null,
	nonce bigint unsigned not null,
	status varchar(16) default 'Generate' not null,
	value varchar(80) not null,
	block_number varchar(80) not null,
	constraint uq_nonce
		unique (from_addr, nonce),
	constraint uq_purchase
		unique (purchase_id)
);

create index idx_tx20_deleted_at
	on tx20 (deleted_at);

create table tx721
(
	id varchar(36) not null
		primary key,
	version bigint default 1 not null,
	created_at bigint not null,
	updated_at bigint not null,
	deleted_at timestamp null,
	purchase_id varchar(36) not null,
	tx_id varchar(66) not null,
	from_addr varchar(42) not null,
	to_addr varchar(42) not null,
	contract varchar(42) not null,
	nonce bigint unsigned not null,
	status varchar(16) default 'Generate' not null,
	token_id bigint not null,
	block_number varchar(80) not null,
	constraint uq_nonce
		unique (from_addr, nonce),
	constraint uq_purchase
		unique (purchase_id)
);

create index idx_tx721_deleted_at
	on tx721 (deleted_at);

