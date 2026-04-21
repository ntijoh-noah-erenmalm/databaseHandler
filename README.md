# databaseHandler
database handler made for exam work. Wont contain every feature of a normal database but will start with a b-tree structure for indexing and data storage




# Structure
- make a b-tree. The b-tree will be after which column of a table is filtered by(int, string, bool, uuid?)
- store b-tree's as bin(?)
- the key has data for pagination, start bit and bit size.
- insert at (first matching deleted data) or at end of pagination
#
- updated rows will update at its current place, if it too large --> append then delete
#
- On delete add the location to free space map.
#
- Command to "vacuum" dead space, defragment space used.
#
- catalog for table column data(types and names)
