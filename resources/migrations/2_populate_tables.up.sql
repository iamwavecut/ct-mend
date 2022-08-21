insert into clients (id, name, code_scan_interval)
values (1, 'Microsoft', 10000),
    (2, 'Apple', 20000),
    (3, 'Alphabet', 5000),
    (4, 'Meta', 1000);

insert into projects (id, client_id, name)
values (1, null, 'Linux'),
    (2, 1, 'Windows 3.11 for Workgroups'),
    (3, 1, 'XBox Fridge Firmware'),
    (4, 2, 'IKettleOs'),
    (6, 2, 'ICar Autopilot Brake on Kernel Panic'),
    (7, 3, 'Brand new messenger to be buried next quarter'),
    (8, 4, 'User data automatic seller');