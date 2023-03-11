CREATE TABLE "users" (
	"id_user" serial NOT NULL,
	"username" varchar(255) NOT NULL UNIQUE,
	"email" varchar(255) NOT NULL UNIQUE,
	"password" varchar(255) NOT NULL,
	CONSTRAINT "voc_user_pk" PRIMARY KEY ("id_user")
) WITH (
  OIDS=FALSE
);

CREATE TABLE "voc_part_sp" (
	"id_part_sp" serial NOT NULL,
	"descrip" varchar(255) NOT NULL UNIQUE,
	CONSTRAINT "voc_part_sp_pk" PRIMARY KEY ("id_part_sp")
) WITH (
  OIDS=FALSE
);

CREATE TABLE "voc_test"(
    "id_test" serial       NOT NULL,
    "descrip" varchar(255) NOT NULL UNIQUE,
    CONSTRAINT "voc_test_pk" PRIMARY KEY ("id_test")
) WITH (
      OIDS= FALSE
    );

CREATE TABLE "words" (
	"part_ps_id" int NOT NULL,
	"word" varchar NOT NULL
) WITH (
  OIDS=FALSE
);

CREATE TABLE "res_test" (
	"user_id" int NOT NULL,
	"test_id" int NOT NULL,
	"date_test" DATE NOT NULL,
	"result_inter" interval,
	"result_1" int,
	"result_2" int
) WITH (
  OIDS=FALSE
);

ALTER TABLE "words" ADD CONSTRAINT "words_fk0" FOREIGN KEY ("part_ps_id") REFERENCES "voc_part_sp"("id_part_sp");
ALTER TABLE "res_test" ADD CONSTRAINT "res_shulte_fk0" FOREIGN KEY ("user_id") REFERENCES "users"("id_user");
ALTER TABLE "res_test" ADD CONSTRAINT "res_shulte_fk1" FOREIGN KEY ("test_id") REFERENCES "voc_test"("id_test");







