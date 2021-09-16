-- qiniu: http(s)://qiniu_public_domain
-- local: api_base_url
-- generic: generic_storage_base_url

UPDATE `user` SET `avatar`=REPLACE(avatar, 'base_url', '');
UPDATE `team` SET `logo`=REPLACE(`logo`, 'base_url', '');
UPDATE `organization` SET `logo`=REPLACE(`logo`, 'base_url', ''), `favicon`=REPLACE(`favicon`, 'base_url', '');