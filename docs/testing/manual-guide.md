# Manual WebSocket Testing Guide

This guide walks you through manually testing the complete chat flow using wscat.

## Setup Information

**User IDs:**
- Alice: `47e4ebd0-2f8c-43f4-b87e-f65cf64e5dcb`
- Bob: `58f56435-97ca-4b02-bc94-247c52807b87`
- Charlie: `aa6779e1-8555-4a17-a4a3-ba3e15ac1975`

**Group ID:**
- Team Alpha: `91d97fba-736e-4d5b-8a22-cb4a75cdb036`

**Tokens:**
- Alice: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU0NzcwNTQsImlhdCI6MTc2NTM5MDY1NCwic3ViIjoiNDdlNGViZDAtMmY4Yy00M2Y0LWI4N2UtZjY1Y2Y2NGU1ZGNiIn0.NzudEFDAGpuXIkarv_OuEztyMpE1S4LI1pvkJbIm1IA`
- Bob: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU0NzcwNTQsImlhdCI6MTc2NTM5MDY1NCwic3ViIjoiNThmNTY0MzUtOTdjYS00YjAyLWJjOTQtMjQ3YzUyODA3Yjg3In0.4jTTeV4nsfNjjRiRgPsV6IAvcx8aQoZUxRqoSRAXp9E`
- Charlie: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU0NzcwNTQsImlhdCI6MTc2NTM5MDY1NCwic3ViIjoiYWE2Nzc5ZTEtODU1NS00YTE3LWE0YTMtYmEzZTE1YWMxOTc1In0.LL__OkZnxSFWU9Dwj_t-MiFG34wUWho5fHwy9XL2dDc`

---

## Testing Steps

### Step 1: Connect Alice via WebSocket

Open a terminal and run:
```bash
wscat -c "ws://localhost:8080/ws?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU0NzcwNTQsImlhdCI6MTc2NTM5MDY1NCwic3ViIjoiNDdlNGViZDAtMmY4Yy00M2Y0LWI4N2UtZjY1Y2Y2NGU1ZGNiIn0.NzudEFDAGpuXIkarv_OuEztyMpE1S4LI1pvkJbIm1IA"
```

You should see: `Connected (press CTRL+C to quit)`

### Step 2: Connect Bob via WebSocket

Open **another terminal** and run:
```bash
wscat -c "ws://localhost:8080/ws?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU0NzcwNTQsImlhdCI6MTc2NTM5MDY1NCwic3ViIjoiNThmNTY0MzUtOTdjYS00YjAyLWJjOTQtMjQ3YzUyODA3Yjg3In0.4jTTeV4nsfNjjRiRgPsV6IAvcx8aQoZUxRqoSRAXp9E"
```

### Step 3: Send DM from Alice to Bob

In Alice's terminal, type:
```json
{"type":"send_message","payload":{"to_user_id":"58f56435-97ca-4b02-bc94-247c52807b87","content":"Hi Bob, this is Alice!"}}
```

**Expected**: Bob's terminal should receive the message immediately.

### Step 4: Send DM from Bob to Alice

In Bob's terminal, type:
```json
{"type":"send_message","payload":{"to_user_id":"47e4ebd0-2f8c-43f4-b87e-f65cf64e5dcb","content":"Hey Alice! How are you?"}}
```

**Expected**: Alice's terminal should receive the message.

### Step 5: Send another DM from Alice to Bob

In Alice's terminal:
```json
{"type":"send_message","payload":{"to_user_id":"58f56435-97ca-4b02-bc94-247c52807b87","content":"I'm good, thanks Bob!"}}
```

### Step 6: Connect Charlie via WebSocket

Open **a third terminal** and run:
```bash
wscat -c "ws://localhost:8080/ws?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU0NzcwNTQsImlhdCI6MTc2NTM5MDY1NCwic3ViIjoiYWE2Nzc5ZTEtODU1NS00YTE3LWE0YTMtYmEzZTE1YWMxOTc1In0.LL__OkZnxSFWU9Dwj_t-MiFG34wUWho5fHwy9XL2dDc"
```

### Step 7: Send Group Message from Alice

In Alice's terminal:
```json
{"type":"send_message","payload":{"group_id":"91d97fba-736e-4d5b-8a22-cb4a75cdb036","content":"Hello Team Alpha!"}}
```

**Expected**: Both Bob and Charlie should receive this message (NOT Alice, since sender doesn't receive own messages).

### Step 8: Send Group Message from Bob

In Bob's terminal:
```json
{"type":"send_message","payload":{"group_id":"91d97fba-736e-4d5b-8a22-cb4a75cdb036","content":"Hey team, Bob here!"}}
```

**Expected**: Alice and Charlie receive it.

### Step 9: Send Group Message from Charlie

In Charlie's terminal:
```json
{"type":"send_message","payload":{"group_id":"91d97fba-736e-4d5b-8a22-cb4a75cdb036","content":"Charlie checking in!"}}
```

**Expected**: Alice and Bob receive it.

---

## Verify Inbox & History APIs

After sending messages via WebSocket, close all wscat connections (CTRL+C) and verify the REST endpoints.

### Check Alice's Conversations (Inbox)

```bash
curl -X GET "http://localhost:8080/conversations" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU0NzcwNTQsImlhdCI6MTc2NTM5MDY1NCwic3ViIjoiNDdlNGViZDAtMmY4Yy00M2Y0LWI4N2UtZjY1Y2Y2NGU1ZGNiIn0.NzudEFDAGpuXIkarv_OuEztyMpE1S4LI1pvkJbIm1IA" | jq
```

**Expected**: Should show 2 conversations (1 DM with Bob, 1 GROUP with Team Alpha).

### Check Bob's Conversations

```bash
curl -X GET "http://localhost:8080/conversations" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU0NzcwNTQsImlhdCI6MTc2NTM5MDY1NCwic3ViIjoiNThmNTY0MzUtOTdjYS00YjAyLWJjOTQtMjQ3YzUyODA3Yjg3In0.4jTTeV4nsfNjjRiRgPsV6IAvcx8aQoZUxRqoSRAXp9E" | jq
```

### Check Charlie's Conversations

```bash
curl -X GET "http://localhost:8080/conversations" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU0NzcwNTQsImlhdCI6MTc2NTM5MDY1NCwic3ViIjoiYWE2Nzc5ZTEtODU1NS00YTE3LWE0YTMtYmEzZTE1YWMxOTc1In0.LL__OkZnxSFWU9Dwj_t-MiFG34wUWho5fHwy9XL2dDc" | jq
```

**Expected**: Should show only 1 conversation (GROUP with Team Alpha), since Charlie hasn't had any DMs.

### Get Alice-Bob DM History

```bash
curl -X GET "http://localhost:8080/messages?target_id=58f56435-97ca-4b02-bc94-247c52807b87&type=DM&limit=10" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU0NzcwNTQsImlhdCI6MTc2NTM5MDY1NCwic3ViIjoiNDdlNGViZDAtMmY4Yy00M2Y0LWI4N2UtZjY1Y2Y2NGU1ZGNiIn0.NzudEFDAGpuXIkarv_OuEztyMpE1S4LI1pvkJbIm1IA" | jq
```

**Expected**: Should show 3 messages between Alice and Bob.

### Get Group Message History

```bash
curl -X GET "http://localhost:8080/messages?target_id=91d97fba-736e-4d5b-8a22-cb4a75cdb036&type=GROUP&limit=20" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU0NzcwNTQsImlhdCI6MTc2NTM5MDY1NCwic3ViIjoiNDdlNGViZDAtMmY4Yy00M2Y0LWI4N2UtZjY1Y2Y2NGU1ZGNiIn0.NzudEFDAGpuXIkarv_OuEztyMpE1S4LI1pvkJbIm1IA" | jq
```

**Expected**: Should show 3 group messages from Alice, Bob, and Charlie.

---

## Checklist

After completing all steps, verify:

- [ ] ✅ Alice and Bob can send DMs to each other
- [ ] ✅ Messages are received in real-time via WebSocket
- [ ] ✅ Group messages broadcast to all members (except sender)
- [ ] ✅ GET /conversations returns correct inbox for all users
- [ ] ✅ Conversations are sorted by last_message_at
- [ ] ✅ GET /messages returns correct DM history
- [ ] ✅ GET /messages returns correct GROUP history
- [ ] ✅ Unread counts are displayed in inbox
- [ ] ✅ Accessing messages resets unread count

---

## Quick Commands Summary

For easy copy-paste, save these as environment variables:

```bash
export ALICE_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU0NzcwNTQsImlhdCI6MTc2NTM5MDY1NCwic3ViIjoiNDdlNGViZDAtMmY4Yy00M2Y0LWI4N2UtZjY1Y2Y2NGU1ZGNiIn0.NzudEFDAGpuXIkarv_OuEztyMpE1S4LI1pvkJbIm1IA"
export BOB_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU0NzcwNTQsImlhdCI6MTc2NTM5MDY1NCwic3ViIjoiNThmNTY0MzUtOTdjYS00YjAyLWJjOTQtMjQ3YzUyODA3Yjg3In0.4jTTeV4nsfNjjRiRgPsV6IAvcx8aQoZUxRqoSRAXp9E"
export CHARLIE_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU0NzcwNTQsImlhdCI6MTc2NTM5MDY1NCwic3ViIjoiYWE2Nzc5ZTEtODU1NS00YTE3LWE0YTMtYmEzZTE1YWMxOTc1In0.LL__OkZnxSFWU9Dwj_t-MiFG34wUWho5fHwy9XL2dDc"
export BOB_ID="58f56435-97ca-4b02-bc94-247c52807b87"
export ALICE_ID="47e4ebd0-2f8c-43f4-b87e-f65cf64e5dcb"
export GROUP_ID="91d97fba-736e-4d5b-8a22-cb4a75cdb036"
```

Then use:
```bash
# Connect Alice
wscat -c "ws://localhost:8080/ws?token=$ALICE_TOKEN"

# Get Alice's inbox
curl -X GET "http://localhost:8080/conversations" -H "Authorization: Bearer $ALICE_TOKEN" | jq

# Get Alice-Bob messages
curl -X GET "http://localhost:8080/messages?target_id=$BOB_ID&type=DM" -H "Authorization: Bearer $ALICE_TOKEN" | jq
```
