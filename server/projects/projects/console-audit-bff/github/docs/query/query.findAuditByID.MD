# Find Audit By ID

The responsibility of this query is to return Audit By ID.

## API

---

### Request

**Method** : `POST`

**URL** : `/query`

**Header** :
- **x-tenant**: Tenant ID
- **x-roles**: Roles

#### Example Header:

| Value    | Example                                 |
|----------|-----------------------------------------|
| x-tenant | TN-9777ef45-ba68-4a52-b4f0-5e85b190cc25 |
| x-roles  | owner                                   |

**Request Body:**

{"query":"query {\n    find_audit_by_id(id: 3644) {\n        id\n        operation {\n            tenant\n            action\n            domain\n            domain_id\n            cid\n            date\n        }\n        user {\n            email\n            roles\n        }\n        user_agent {\n            device\n            device_ip\n        }\n        http {\n            code\n            request\n            response\n        }\n    }\n}","variables":{}}

---

### Response

---

{
    "data": {
        "find_audit_by_id": {
            "id": 3644,
            "operation": {
                "tenant": "08102A0B-D4F8-42A2-8B0E-2052D05577D7",
                "action": "create",
                "domain": "user",
                "domain_id": "leonardo.carmona+pismo@pismo.io",
                "cid": "e68394eb-858e-9508-9078-32e17857c5aa",
                "date": "2022-09-19 16:59:15"
            },
            "user": {
                "email": "pismo@pismo.com",
                "roles": [
                    "account-holder",
                    "backoffice-operator",
                    "crm-operator",
                    "dashboards-operator",
                    "external-operator",
                    "marketplace-administrator",
                    "marketplace-manager",
                    "marketplace-merchant-manager",
                    "owner",
                    "parameters-manager",
                    "parameters-operator",
                    "reports-files-accounting",
                    "reports-files-financial",
                    "reports-files-fiscal",
                    "reports-files-management",
                    "reports-files-regulatory",
                    "tokpag"
                ]
            },
            "user_agent": {
                "device": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
                "device_ip": "2804:1b0:1603:4a0f:3d67:ae3b:a2dc:841f"
            },
            "http": {
                "code": 201,
                "request": "",
                "response": ""
            }
        }
    }
}

#### Success

**Code** : `200 - OK`

**URL** : `POST /query`

**Response body:**
{
    "data": {
        "find_audit_by_id": {
            "id": 3644,
            "operation": {
                "tenant": "08102A0B-D4F8-42A2-8B0E-2052D05577D7",
                "action": "create",
                "domain": "user",
                "domain_id": "leonardo.carmona+pismo@pismo.io",
                "cid": "e68394eb-858e-9508-9078-32e17857c5aa",
                "date": "2022-09-19 16:59:15"
            },
            "user": {
                "email": "pismo@pismo.com",
                "roles": [
                    "account-holder",
                    "backoffice-operator",
                    "crm-operator",
                    "dashboards-operator",
                    "external-operator",
                    "marketplace-administrator",
                    "marketplace-manager",
                    "marketplace-merchant-manager",
                    "owner",
                    "parameters-manager",
                    "parameters-operator",
                    "reports-files-accounting",
                    "reports-files-financial",
                    "reports-files-fiscal",
                    "reports-files-management",
                    "reports-files-regulatory",
                    "tokpag"
                ]
            },
            "user_agent": {
                "device": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
                "device_ip": "2804:1b0:1603:4a0f:3d67:ae3b:a2dc:841f"
            },
            "http": {
                "code": 201,
                "request": "",
                "response": ""
            }
        }
    }
}

---

**Code** : `200 - Not Found`

**URL** : `POST /query`

**Response body:**
{
    "data": {
        "find_audit_by_id": null
    }
}
---

#### Error

**Code** : `500 - Internal Server Error`

**URL** : `POST /query`

**Response body:**
{
    "message": "internal_server_error",
    "extensions": {
        "tenant": "TN-xxx",
        "user": {
            "email": "pismo@pismo.io",
            "roles": [
                "owner",
            ]
        },
        "tracking": {
            "tracking_id": "4c30b2651af34dcaeab2e144339d0179",
            "correlation_id": "50bb0ec8-53c4-443d-abdf-765d772adbfa"
        },
        "params": {}
    }
}

---

## Sequence Diagram

<div align="center">
  <img alt="findAuditByID" src="query.findAuditByID.png"/>
</div>
