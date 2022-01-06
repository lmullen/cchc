ALTER TABLE files
  DROP CONSTRAINT files_item_id_fkey,
  ADD CONSTRAINT files_item_id_fkey FOREIGN KEY (item_id) REFERENCES items (id);

ALTER TABLE items_in_collections
  DROP CONSTRAINT items_in_collections_item_id_fkey,
  ADD CONSTRAINT items_in_collections_item_id_fkey FOREIGN KEY (item_id) REFERENCES items (id);

ALTER TABLE resources
  DROP CONSTRAINT resources_item_id_fkey,
  ADD CONSTRAINT resources_item_id_fkey FOREIGN KEY (item_id) REFERENCES items (id);

