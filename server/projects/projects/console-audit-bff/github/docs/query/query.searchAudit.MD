# Search Audit

The responsibility of this query is to return a paginated audit by parameters.

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

{"query":"query {\n    search_audit(per_page: 1, page: 1, order: asc, begin_date: \"2022-08-01T00:00:00Z\", end_date: \"2022-09-01T00:00:00Z\", user: \"pismo@pismo.io\", actions: [create], domains: [\"user\"], domain_id: \"CID-crm-operator\", cid: \"xxx\") {\n        current_page\n        pages\n        per_page\n        total_items\n        audits {\n            id\n        }\n    }\n}","variables":{}}

---

### Response

---

{
    "data": {
        "search_audit": {
            "current_page": 1,
            "page": 1,
            "per_page": 1,
            "total_items": 1,
            "audits": [
                {
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
            ]
        }
    }
}

#### Success

**Code** : `200 - OK`

**URL** : `POST /query`

**Response body:**
{
    "data": {
        "search_audit": {
            "current_page": 1,
            "page": 1,
            "per_page": 1,
            "total_items": 1,
            "audits": [
                {
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
            ]
        }
    }
}

---

#### Error

**Code** : `500 - Internal Server Error`

**URL** : `POST /query`

**Response body:**
{
    "errors": [
        {
            "message": "internal_server_error",
            "path": [
                "search_audit"
            ],
            "extensions": {
                "error": {
                    "tenant": "08102A0B-D4F8-42A2-8B0E-2052D05577D7",
                    "user": {
                        "email": "",
                        "roles": [
                            "owner"
                        ]
                    },
                    "tracking": {
                        "tracking_id": "1afbfadf048f6109b3671af23ae9e17c",
                        "correlation_id": "10b78a70-12aa-4549-b266-16129fb3d232"
                    },
                    "params": {}
                }
            }
        }
    ],
    "data": null
}

---

## Sequence Diagram

<div align="center">
  <img alt="searchAudit" src="query.searchAudit.png"/>
</div>
