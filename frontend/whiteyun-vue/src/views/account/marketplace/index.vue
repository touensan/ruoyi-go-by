<template>
  <div class="app-container">
    <section class="market-hero">
      <span>用户服务入口</span>
      <h1>服务市场</h1>
      <p>这是 RuoYi-Go BY 为业务项目预留的用户侧商品入口。项目可以在这里接入 API 套餐、应用授权、订阅或虚拟商品，并统一使用平台积分结算。</p>
      <div class="market-actions"><el-button type="primary" @click="$router.push('/account/points')">查看我的积分</el-button><el-button @click="$router.push('/account/exchange')">前往兑换中心</el-button></div>
    </section>
  </div>
</template>
<script setup lang="ts" name="AccountMarketplace" />
<style scoped lang="scss">
.market-hero { max-width: 900px; padding: 34px; border: 1px solid var(--admin-line); border-radius: 16px; background: var(--admin-surface); }
.market-hero span { color: var(--el-color-primary); font-size: 13px; font-weight: 700; }.market-hero h1 { margin: 8px 0 14px; font-size: 42px; }
.market-hero p { max-width: 720px; color: var(--el-text-color-secondary); font-size: 16px; line-height: 1.8; }.market-actions { display: flex; gap: 10px; margin-top: 24px; }
</style>
