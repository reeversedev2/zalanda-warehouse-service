# zalanda-warehouse-service

This microservice handles the Products and Company related data. The service is directly relying on PostgresDB for database.

## How a warehouse operates (fast mental model) ￼
 • Inbound/Receiving: Goods arrive against a purchase order (PO) or without one. Staff verify quantities/condition, record what actually arrived, and create a receipt. Items often get labeled (SKU/lot/serial/barcode). Discrepancies (over/short/damaged) are captured.
 • Putaway: Received items are moved from the dock to storage bins. Either the system suggests bins (directed putaway) or the operator picks one. Inventory is now at-bin accurate.
 • Storage & Inventory: Stock is tracked by SKU at specific bins (and sometimes lot/batch, serial, and unit of measure). All stock changes go through controlled transactions to keep “book” inventory correct.
 • Order fulfillment (picking): Customer or internal orders reserve stock. Work is released as picklists. Pickers go to bins, confirm picks, and reduce inventory from those bins.
 • Packing & Shipping: Picked items are packed into cartons. Shipping info (carrier, service, tracking) is recorded; labels may be printed. The order is shipped and inventory is decremented/finalized.
 • Inventory control: Cycle counts verify bins regularly; adjustments are made with reasons to keep accuracy high.
 • Returns (later): Returned goods are received with a disposition (good, rework, scrap).
