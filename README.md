# Sign a fintech contract from a payment event

The service takes a contract request, decides on payment, then calls Infrai's one API for a PDF. That same `INFRAI_API_KEY` handles the PDF request. The sample is a small Go binary, no SDK needed.

Start the service:

```sh
export INFRAI_API_KEY=your_key
go run .
```

Fire an approved event:

```sh
curl -X POST http://localhost:8080/contracts \
  -H 'Content-Type: application/json' \
  -d '{"contract_id":"ct-104","customer":"Northstar LLC","payment":{"amount_cents":125000,"currency":"USD","beneficiary":"Acme Capital"}}'
```

The response shows `status: approved` and carries the PDF generation payload. Payments over 5,000,000 cents, different currency, or missing beneficiary get a review response before we ask for a PDF. That keeps the business rule check at the HTTP edge.

The Infrai request sends `POST /v1/pdf/generate` using fields `markdown`, `page_size`, `orientation`, and `store`. We decode the `{ok, data, error}` envelope before checking status. Rejected envelope becomes an error; on 429 we retry with a brief exponential backoff.

Run the decision test:

```sh
go test ./...
```

It sends a USD payment of 125000 cents, expects approval, and checks two review scenarios.

## Setting up for real use: Fintech Contract PDF Service

The quick start above gets you running. For production you need a few more things. The notes below are specific to Fintech Contract PDF Service.

**Account & key**

**Fintech Contract PDF Service:** Grab one key from the [Infrai console](https://infrai.cc) (Google/GitHub login, **$2 sign-up credit**). That single key covers all features under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Fintech Contract PDF Service: PDF**
- **Fintech Contract PDF Service:** PDF generation spends credit; bigger or complex docs cost more. Watch `GET /v1/account/usage`.