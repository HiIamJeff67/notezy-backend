## Created Product

- Request Body:

```JSON
{
  "name": "Upgrade Plan Product",
  "description": "Product for upgraded subscription",
  "type": "SERVICE",
  "category": "SOFTWARE"
}
```

- Response Body:

```JSON
{
    "id": "PROD-22007905DJ014973Y",
    "name": "Upgrade Plan Product",
    "description": "Product for upgraded subscription",
    "create_time": "2026-04-06T09:16:15Z",
    "links": [
        {
            "href": "https://api.sandbox.paypal.com/v1/catalogs/products/PROD-22007905DJ014973Y",
            "rel": "self",
            "method": "GET"
        },
        {
            "href": "https://api.sandbox.paypal.com/v1/catalogs/products/PROD-22007905DJ014973Y",
            "rel": "edit",
            "method": "PATCH"
        }
    ]
}
```

---

## Created Monthly Pro Plan

- **Request** Body:

```JSON
{
    "product_id": "PROD-22007905DJ014973Y",
    "name": "Notezy Monthly Pro Plan",
    "description": "Upgrade the user plan of Notezy to monthly pro plan",
    "status": "ACTIVE",
    "billing_cycles": [
        {
            "frequency": {
                "interval_unit": "MONTH",
                "interval_count": 1
            },
            "tenure_type": "TRIAL",
            "sequence": 1,
            "total_cycles": 1,
            "pricing_scheme": {
                "fixed_price": {
                    "value": "2.49",
                    "currency_code": "USD"
                }
            }
        },
        {
            "frequency": {
                "interval_unit": "MONTH",
                "interval_count": 1
            },
            "tenure_type": "REGULAR",
            "sequence": 2,
            "total_cycles": 0,
            "pricing_scheme": {
                "fixed_price": {
                    "value": "4.99",
                    "currency_code": "USD"
                }
            }
        }
  ],
  "payment_preferences": {
    "auto_bill_outstanding": true,
    "setup_fee_failure_action": "CONTINUE",
    "payment_failure_threshold": 3
  }
}
```

- Response Body:

```JSON
{
    "id": "P-4LN51972TD528344JNHJX5NQ",
    "product_id": "PROD-22007905DJ014973Y",
    "name": "Notezy Monthly Pro Plan",
    "status": "ACTIVE",
    "description": "Upgrade the user plan of Notezy to monthly pro plan",
    "usage_type": "LICENSED",
    "create_time": "2026-04-06T09:36:54Z",
    "links": [
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-4LN51972TD528344JNHJX5NQ",
            "rel": "self",
            "method": "GET",
            "encType": "application/json"
        },
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-4LN51972TD528344JNHJX5NQ",
            "rel": "edit",
            "method": "PATCH",
            "encType": "application/json"
        },
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-4LN51972TD528344JNHJX5NQ/deactivate",
            "rel": "self",
            "method": "POST",
            "encType": "application/json"
        }
    ]
}
```

---

## Created Yearly Pro Plan

- **Request** Body:

```JSON
{
  "product_id": "PROD-22007905DJ014973Y",
  "name": "Notezy Yearly Pro Plan",
  "description": "Upgrade the user plan of Notezy to yearly pro plan",
  "status": "ACTIVE",
  "billing_cycles": [
    {
      "frequency": {
        "interval_unit": "YEAR",
        "interval_count": 1
      },
      "tenure_type": "REGULAR",
      "sequence": 1,
      "total_cycles": 0,
      "pricing_scheme": {
        "fixed_price": {
          "value": "49.99",
          "currency_code": "USD"
        }
      }
    }
  ],
  "payment_preferences": {
    "auto_bill_outstanding": true,
    "setup_fee_failure_action": "CONTINUE",
    "payment_failure_threshold": 3
  }
}
```

- Response Body

```JSON
{
    "id": "P-9MB559415V1980509NHJX7LY",
    "product_id": "PROD-22007905DJ014973Y",
    "name": "Notezy Yearly Pro Plan",
    "status": "ACTIVE",
    "description": "Upgrade the user plan of Notezy to yearly pro plan",
    "usage_type": "LICENSED",
    "create_time": "2026-04-06T09:41:03Z",
    "links": [
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-9MB559415V1980509NHJX7LY",
            "rel": "self",
            "method": "GET",
            "encType": "application/json"
        },
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-9MB559415V1980509NHJX7LY",
            "rel": "edit",
            "method": "PATCH",
            "encType": "application/json"
        },
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-9MB559415V1980509NHJX7LY/deactivate",
            "rel": "self",
            "method": "POST",
            "encType": "application/json"
        }
    ]
}
```

---

## Created Monthly Premium Plan

- **Request** Body:

```JSON
{
  "product_id": "PROD-22007905DJ014973Y",
  "name": "Notezy Monthly Premium Plan",
  "description": "Upgrade the user plan of Notezy to monthly premium plan",
  "status": "ACTIVE",
  "billing_cycles": [
    {
      "frequency": {
        "interval_unit": "MONTH",
        "interval_count": 1
      },
      "tenure_type": "REGULAR",
      "sequence": 1,
      "total_cycles": 0,
      "pricing_scheme": {
        "fixed_price": {
          "value": "9.99",
          "currency_code": "USD"
        }
      }
    }
  ],
  "payment_preferences": {
    "auto_bill_outstanding": true,
    "setup_fee_failure_action": "CONTINUE",
    "payment_failure_threshold": 3
  }
}
```

- Response Body:

```JSON
{
    "id": "P-351611974A6912332NHJYECY",
    "product_id": "PROD-22007905DJ014973Y",
    "name": "Notezy Monthly Premium Plan",
    "status": "ACTIVE",
    "description": "Upgrade the user plan of Notezy to monthly premium plan",
    "usage_type": "LICENSED",
    "create_time": "2026-04-06T09:51:07Z",
    "links": [
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-351611974A6912332NHJYECY",
            "rel": "self",
            "method": "GET",
            "encType": "application/json"
        },
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-351611974A6912332NHJYECY",
            "rel": "edit",
            "method": "PATCH",
            "encType": "application/json"
        },
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-351611974A6912332NHJYECY/deactivate",
            "rel": "self",
            "method": "POST",
            "encType": "application/json"
        }
    ]
}
```

---

## Created Yearly Premium Plan

- **Request** Body:

```JSON
{
  "product_id": "PROD-22007905DJ014973Y",
  "name": "Notezy Yearly Premium Plan",
  "description": "Upgrade the user plan of Notezy to yearly premium plan",
  "status": "ACTIVE",
  "billing_cycles": [
    {
      "frequency": {
        "interval_unit": "YEAR",
        "interval_count": 1
      },
      "tenure_type": "REGULAR",
      "sequence": 1,
      "total_cycles": 0,
      "pricing_scheme": {
        "fixed_price": {
          "value": "99.99",
          "currency_code": "USD"
        }
      }
    }
  ],
  "payment_preferences": {
    "auto_bill_outstanding": true,
    "setup_fee_failure_action": "CONTINUE",
    "payment_failure_threshold": 3
  }
}
```

- Response Body:

```JSON
{
    "id": "P-84627481GN3337838NHJYEJA",
    "product_id": "PROD-22007905DJ014973Y",
    "name": "Notezy Yearly Premium Plan",
    "status": "ACTIVE",
    "description": "Upgrade the user plan of Notezy to yearly premium plan",
    "usage_type": "LICENSED",
    "create_time": "2026-04-06T09:51:32Z",
    "links": [
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-84627481GN3337838NHJYEJA",
            "rel": "self",
            "method": "GET",
            "encType": "application/json"
        },
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-84627481GN3337838NHJYEJA",
            "rel": "edit",
            "method": "PATCH",
            "encType": "application/json"
        },
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-84627481GN3337838NHJYEJA/deactivate",
            "rel": "self",
            "method": "POST",
            "encType": "application/json"
        }
    ]
}
```

---

## Created Monthly Ultimate Plan

- **Request** Body:

```JSON
{
  "product_id": "PROD-22007905DJ014973Y",
  "name": "Notezy Monthly Ultimate Plan",
  "description": "Upgrade the user plan of Notezy to monthly ultimate plan",
  "status": "ACTIVE",
  "billing_cycles": [
    {
      "frequency": {
        "interval_unit": "MONTH",
        "interval_count": 1
      },
      "tenure_type": "REGULAR",
      "sequence": 1,
      "total_cycles": 0,
      "pricing_scheme": {
        "fixed_price": {
          "value": "19.99",
          "currency_code": "USD"
        }
      }
    }
  ],
  "payment_preferences": {
    "auto_bill_outstanding": true,
    "setup_fee_failure_action": "CONTINUE",
    "payment_failure_threshold": 3
  }
}
```

- Response Body:

```JSON
{
    "id": "P-3B912255TH6394814NHJYEMA",
    "product_id": "PROD-22007905DJ014973Y",
    "name": "Notezy Monthly Ultimate Plan",
    "status": "ACTIVE",
    "description": "Upgrade the user plan of Notezy to monthly ultimate plan",
    "usage_type": "LICENSED",
    "create_time": "2026-04-06T09:51:44Z",
    "links": [
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-3B912255TH6394814NHJYEMA",
            "rel": "self",
            "method": "GET",
            "encType": "application/json"
        },
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-3B912255TH6394814NHJYEMA",
            "rel": "edit",
            "method": "PATCH",
            "encType": "application/json"
        },
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-3B912255TH6394814NHJYEMA/deactivate",
            "rel": "self",
            "method": "POST",
            "encType": "application/json"
        }
    ]
}
```

---

## Created Yearly Ultimate Plan

- **Request** Body:

```JSON
{
  "product_id": "PROD-22007905DJ014973Y",
  "name": "Notezy Yearly Ultimate Plan",
  "description": "Upgrade the user plan of Notezy to yearly ultimate plan",
  "status": "ACTIVE",
  "billing_cycles": [
    {
      "frequency": {
        "interval_unit": "YEAR",
        "interval_count": 1
      },
      "tenure_type": "REGULAR",
      "sequence": 1,
      "total_cycles": 0,
      "pricing_scheme": {
        "fixed_price": {
          "value": "199.99",
          "currency_code": "USD"
        }
      }
    }
  ],
  "payment_preferences": {
    "auto_bill_outstanding": true,
    "setup_fee_failure_action": "CONTINUE",
    "payment_failure_threshold": 3
  }
}
```

- Response Body:

```JSON
{
    "id": "P-4WS50500MM359840MNHJYEOY",
    "product_id": "PROD-22007905DJ014973Y",
    "name": "Notezy Yearly Ultimate Plan",
    "status": "ACTIVE",
    "description": "Upgrade the user plan of Notezy to yearly ultimate plan",
    "usage_type": "LICENSED",
    "create_time": "2026-04-06T09:51:55Z",
    "links": [
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-4WS50500MM359840MNHJYEOY",
            "rel": "self",
            "method": "GET",
            "encType": "application/json"
        },
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-4WS50500MM359840MNHJYEOY",
            "rel": "edit",
            "method": "PATCH",
            "encType": "application/json"
        },
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-4WS50500MM359840MNHJYEOY/deactivate",
            "rel": "self",
            "method": "POST",
            "encType": "application/json"
        }
    ]
}
```

---

## Created Monthly Enterprise Plan

- **Request** Body:

```JSON
{
  "product_id": "PROD-22007905DJ014973Y",
  "name": "Notezy Monthly Enterprise Plan",
  "description": "Upgrade the user plan of Notezy to monthly enterprise plan",
  "status": "ACTIVE",
  "billing_cycles": [
    {
      "frequency": {
        "interval_unit": "MONTH",
        "interval_count": 1
      },
      "tenure_type": "REGULAR",
      "sequence": 1,
      "total_cycles": 0,
      "pricing_scheme": {
        "fixed_price": {
          "value": "49.99",
          "currency_code": "USD"
        }
      }
    }
  ],
  "payment_preferences": {
    "auto_bill_outstanding": true,
    "setup_fee_failure_action": "CONTINUE",
    "payment_failure_threshold": 3
  }
}
```

- Response Body:

```JSON
{
    "id": "P-9XP882067J683411BNHJYERI",
    "product_id": "PROD-22007905DJ014973Y",
    "name": "Notezy Monthly Enterprise Plan",
    "status": "ACTIVE",
    "description": "Upgrade the user plan of Notezy to monthly enterprise plan",
    "usage_type": "LICENSED",
    "create_time": "2026-04-06T09:52:05Z",
    "links": [
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-9XP882067J683411BNHJYERI",
            "rel": "self",
            "method": "GET",
            "encType": "application/json"
        },
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-9XP882067J683411BNHJYERI",
            "rel": "edit",
            "method": "PATCH",
            "encType": "application/json"
        },
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-9XP882067J683411BNHJYERI/deactivate",
            "rel": "self",
            "method": "POST",
            "encType": "application/json"
        }
    ]
}
```

---

## Created Yearly Enterprise Plan

- **Request** Body:

```JSON
{
  "product_id": "PROD-22007905DJ014973Y",
  "name": "Notezy Yearly Enterprise Plan",
  "description": "Upgrade the user plan of Notezy to yearly enterprise plan",
  "status": "ACTIVE",
  "billing_cycles": [
    {
      "frequency": {
        "interval_unit": "YEAR",
        "interval_count": 1
      },
      "tenure_type": "REGULAR",
      "sequence": 1,
      "total_cycles": 0,
      "pricing_scheme": {
        "fixed_price": {
          "value": "499.99",
          "currency_code": "USD"
        }
      }
    }
  ],
  "payment_preferences": {
    "auto_bill_outstanding": true,
    "setup_fee_failure_action": "CONTINUE",
    "payment_failure_threshold": 3
  }
}
```

- Response Body:

```JSON
{
    "id": "P-2PT73314S5217944VNHJYEYI",
    "product_id": "PROD-22007905DJ014973Y",
    "name": "Notezy Yearly Enterprise Plan",
    "status": "ACTIVE",
    "description": "Upgrade the user plan of Notezy to yearly enterprise plan",
    "usage_type": "LICENSED",
    "create_time": "2026-04-06T09:52:33Z",
    "links": [
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-2PT73314S5217944VNHJYEYI",
            "rel": "self",
            "method": "GET",
            "encType": "application/json"
        },
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-2PT73314S5217944VNHJYEYI",
            "rel": "edit",
            "method": "PATCH",
            "encType": "application/json"
        },
        {
            "href": "https://api.sandbox.paypal.com/v1/billing/plans/P-2PT73314S5217944VNHJYEYI/deactivate",
            "rel": "self",
            "method": "POST",
            "encType": "application/json"
        }
    ]
}
```
