#!/usr/bin/env python3
"""Verify the v1.4.0 points/exchange foundation invariants without a live database."""

from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
sql = (ROOT / "sql/upgrade/v1.4.0-points-exchange.sql").read_text(encoding="utf-8")
service = (ROOT / "app/service/platform_service.go").read_text(encoding="utf-8")
sidebar = (ROOT / "frontend/whiteyun-vue/src/layout/components/Sidebar/index.vue").read_text(encoding="utf-8")
package = (ROOT / "frontend/whiteyun-vue/package.json").read_text(encoding="utf-8")
config = (ROOT / "application-example.yaml").read_text(encoding="utf-8")

checks = {
    "point account schema": "CREATE TABLE IF NOT EXISTS platform_point_accounts" in sql,
    "recharge schema": "CREATE TABLE IF NOT EXISTS platform_point_recharge_orders" in sql,
    "task and code schema": "platform_reward_tasks" in sql and "platform_exchange_codes" in sql,
    "rainyun default off": '"rainyunEnabled":false' in sql,
    "code plaintext not stored": "code_hmac BINARY(32)" in sql and "code_value" not in sql,
    "one claim per account/task": "uk_platform_reward_claim_account_task" in sql,
    "one provider subject": "uk_platform_reward_claim_subject" in sql,
    "row locks used": 'clause.Locking{Strength: "UPDATE"}' in service,
    "hmac code digest": "hmac.New(sha256.New" in service,
    "account divider": "account-section-divider" in sidebar and "我的服务" in sidebar,
    "version synchronized": '"version": "1.4.0"' in package and "version: 1.4.0" in config,
}

failed = [name for name, passed in checks.items() if not passed]
for name, passed in checks.items():
    print(f"[{'PASS' if passed else 'FAIL'}] {name}")
if failed:
    raise SystemExit("v1.4.0 verification failed: " + ", ".join(failed))
