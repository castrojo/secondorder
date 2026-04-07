# Tech Spec: Strategic Alignment Data Model & API

## Overview
This change introduces the `ApexBlock` data model and updates the `WorkBlock` model to support strategic alignment and North Star metrics.

## Data Model Changes

### ApexBlock
A new top-level strategic block.
- `id`: UUID
- `title`: String
- `goal`: String
- `status`: String (active, archived)
- `created_at`: DateTime
- `updated_at`: DateTime

### WorkBlock (Updated)
Added fields for strategic alignment:
- `north_star_metric`: String
- `north_star_target`: String
- `apex_block_id`: Foreign key to `apex_blocks.id`

## API Changes

### Apex Blocks
- `GET /api/v1/apex-blocks`: List all apex blocks.
- `GET /api/v1/apex-blocks/{id}`: Get details of a specific apex block.
- `POST /api/v1/apex-blocks`: Create a new apex block.
- `PATCH /api/v1/apex-blocks/{id}`: Update an existing apex block.

### Work Blocks
- `GET /api/v1/work-blocks`: Updated to include new fields.
- `POST /api/v1/work-blocks`: Updated to support creating with North Star metrics and Apex Block association.
- `PATCH /api/v1/work-blocks/{id}`: Updated to support updating all fields including status and strategic alignment fields.

## Database Migrations
- `013_apex_blocks.sql`: Creates `apex_blocks` table and alters `work_blocks` table.
