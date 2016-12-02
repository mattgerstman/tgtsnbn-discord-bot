CREATE TABLE users (user_id varchar(128), guild_id varchar(128), house varchar(32), num_points int);
ALTER TABLE users ADD PRIMARY KEY (user_id, guild_id);
