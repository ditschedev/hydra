-- Migration generated by the command below; DO NOT EDIT.
-- hydra:generate hydra migrations gen

ALTER TABLE hydra_oauth2_code DROP FOREIGN KEY hydra_oauth2_code_challenge_id_fk;
