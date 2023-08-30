CREATE TABLE users(
                      id serial not null unique,
                      username varchar(255) not null unique
);

CREATE TABLE segments(
                    id serial not null unique,
                    segment varchar(255) not null unique
);


CREATE TABLE segmentation_info (
                                     id serial not null unique,
                                     user_id int REFERENCES users(id) ON DELETE CASCADE,
                                     segments_list varchar(255)
);

CREATE TABLE segmentation_history (
                                     id serial not null unique,
                                     user_id int REFERENCES users(id) ON DELETE CASCADE,
                                     segment varchar(255),
                                     status varchar(255) not null,
                                     date timestamp
);
