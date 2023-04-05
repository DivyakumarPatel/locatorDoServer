

create table my_users(
    id INT GENERATED ALWAYS AS IDENTITY,
    email varchar(50) unique  not null,
    date_created timestamp default CURRENT_TIMESTAMP,
    first_name varchar(50),
    middle_name varchar(50) ,
    phone_number varchar(50),
    password text not null,
    firebase_id text,
    is_admin bool default false,
    profile_photo text default "https://www.pngitem.com/pimgs/m/30-307416_profile-icon-png-image-free-download-searchpng-employee.png",


    primary key (id)
)

ALTER TABLE users
ADD COLUMN profile_photo TEXT DEFAULT 'https://www.pngitem.com/pimgs/m/30-307416_profile-icon-png-image-free-download-searchpng-employee.png';




create table distance(
    user_id int,
    current_latitude double precision,
    current_longitude double precision,
    max_distance double precision,
    current_latitude double precision,
    current_longitude double precision,
    latest_update timestamp default CURRENT_TIMESTAMP,

    primary key (user_id),

    CONSTRAINT fk_users
        FOREIGN KEY(user_id)
            REFERENCES product(id)
)
user1
test2
test3
new
fahim
test
charles
Hussnain