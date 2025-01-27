# This script imports DataStore data from the Pokken Tournament standalone server, github.com/PretendoNetwork/pokken-tournament
# This old data was stored as files in directories, but now we use a PostgreSQL database.

# Usage: python scripts/import_legacy_pokken_data.py path/to/old/data
# Please run meganex with the Pokken config file at least once to set up the destination database schema.

from datetime import datetime
import sys
import dotenv
import psycopg
import pathlib

pguri = dotenv.dotenv_values()["PN_MEGANEX_POSTGRES_URI"]

items: list[dict] = []

data_dir = pathlib.Path(sys.argv[1]) / "0"
for meta in data_dir.glob("*.meta"):
    data = meta.with_name(meta.stem)
    if not meta.exists() or not data.exists():
        print(f"Skipping bad {meta}")
        continue
    with open(meta) as f:
        metadata = [s.strip() for s in f.readlines()]
    if len(metadata) != 9:
        print(f"Skipping bad {meta}")
        continue

    with open(data, 'rb') as f:
        metabinary = f.read()

    item = {
        "owner": data.stem,
        "size": int(metadata[0]),
        "name": metadata[1],
        "data_type": int(metadata[2]),
        "meta_binary": metabinary,
        "permission": metadata[3].split(";")[0],
        "permission_recipients": [int(r) for r in metadata[3].split(";")[1].split(",") if r],
        "delete_permission": metadata[4].split(";")[0],
        "delete_permission_recipients": [int(r) for r in metadata[4].split(";")[1].split(",") if r],
        "creation_date": datetime.fromisoformat(metadata[5]),
        "update_date": datetime.fromisoformat(metadata[6]),
        "period": int(metadata[7]),
        "flag": int(metadata[8]),
        "refer_data_id": 0,
        "tags": [],
        "persistence_slot_id": 0,
        "extra_data": [],
    }

    items.append(item)

print(f"[ OK ] Read {len(items)} items.")
ok = 0
with psycopg.connect(pguri) as conn:
    with conn.cursor() as cur:
        for item in items:
            cur.execute("""
            INSERT INTO datastore.objects (owner, size, name, data_type, meta_binary, permission, permission_recipients, delete_permission, delete_permission_recipients, flag, period, refer_data_id, tags, persistence_slot_id, extra_data, creation_date, update_date)
            VALUES (%(owner)s, %(size)s, %(name)s, %(data_type)s, %(meta_binary)s, %(permission)s, %(permission_recipients)s, %(delete_permission)s, %(delete_permission_recipients)s, %(flag)s, %(period)s, %(refer_data_id)s, %(tags)s, %(persistence_slot_id)s, %(extra_data)s, %(creation_date)s, %(update_date)s);
            """, item)
            ok += 1
            print(f"[ OK ] Inserted {item['owner']} - {ok}/{len(items)}")