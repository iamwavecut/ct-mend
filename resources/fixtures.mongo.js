db.clients.insertMany(
    [
        {"id": 1, "name": "Microsoft", "settings": {"code_scan_interval": 10000}},
        {
            "id": 2,
            "name": "Apple",
            "settings": {"code_scan_interval": 20000}
        }, {"id": 3, "name": "Alphabet", "settings": {"code_scan_interval": 5000}},
        {
            "id": 4,
            "name": "Meta",
            "settings": {"code_scan_interval": 1000}
        }
    ]
)
db.projects.insertMany(
    [
        {"id": 1, "name": "Linux"},
        {"id": 2, "client_id": 1, "name": "Windows 3.11 for Workgroups"},
        {"id": 3, "client_id": 1, "name": "XBox Fridge Firmware"},
        {"id": 4, "client_id": 2, "name": "IKettleOs"},
        {"id": 6, "client_id": 2, "name": "ICar Autopilot Brake on Kernel Panic"},
        {"id": 7, "client_id": 3, "name": "Brand new messenger to be buried next quarter"},
        {"id": 8, "client_id": 4, "name": "User data automatic seller"}
    ]
)
counters = [
    {"_id": "clients", seq: 5},
    {"_id": "projects", seq: 9},
];
db.counters.insertMany(counters)