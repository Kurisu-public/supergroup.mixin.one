<template>
  <loading :loading="loading" :fullscreen="true">
    <div class="prepare-packet-page">
      <nav-bar :title="$t('prepare_packet.title')" :hasTopRight="false" :hasBack="true"></nav-bar>
      <van-cell-group title="">
        <row-select
          :index="0"
          :title="$t('prepare_packet.select_assets')"
          :columns="assets"
          placeholder="Tap to Select"
          @change="onChangeAsset">
          <span slot="text">{{selectedAsset ? selectedAsset.text : 'Tap to Select'}}</span>
        </row-select>
        <van-cell>
          <van-field type="number" v-model="form.amount" :label="$t('prepare_packet.amount')" :placeholder="$t('prepare_packet.placeholder_amount')">
            <span slot="right-icon">{{selectedAsset ? selectedAsset.symbol : ''}}</span>
          </van-field>
        </van-cell>
        <van-cell>
          <van-field type="number" v-model="form.shares" :label="$t('prepare_packet.shares')" :placeholder="$t('prepare_packet.placeholder_shares', {count: participantsCount})">
          </van-field>
        </van-cell>
        <van-cell>
          <van-field v-model="form.memo" :label="$t('prepare_packet.memo')" :placeholder="$t('prepare_packet.placeholder_memo')">
          </van-field>
        </van-cell>
      </van-cell-group>
      <van-row style="padding: 20px">
        <van-col span="24">
          <van-button style="width: 100%" type="info" :disabled="!validated" @click="pay">{{$t('prepare_packet.pay')}}</van-button>
        </van-col>
      </van-row>

    </div>
  </loading>
</template>

<script>
import NavBar from '@/components/Nav'
import RowSelect from '@/components/RowSelect'
import Row from '@/components/Nav'
import Loading from '@/components/Loading'
import uuid from 'uuid'
import {Toast} from 'vant'
import { CLIENT_ID } from '@/constants'

export default {
  name: 'Prepare-Packet',
  props: {
    msg: String
  },
  data () {
    return {
      loading: false,
      coversationId: '',
      participantsCount: 0,
      assets: [],
      selectedAsset: null,
      form: {
        amount: '',
        shares: '',
        memo: this.$t('prepare_packet.default_memo', {symbol: 'BTC'})
      }
    }
  },
  components: {
    NavBar, RowSelect, Loading
  },
  async mounted () {
    const packetMaxLimit = 200
    this.loading = true
    let prepareInfo = await this.GLOBAL.api.packet.prepare()
    if (prepareInfo) {
      this.assets = prepareInfo.data.assets.map((x) => {
        x.text = `${x.symbol} (${x.balance})`
        return x
      })
      if (this.assets.length) {
        this.selectedAsset = this.assets[0]
        this.form.memo = this.$t('prepare_packet.default_memo', {symbol: this.selectedAsset.symbol})
      }
      this.coversationId = prepareInfo.data.conversation.coversation_id

      if (prepareInfo.data.conversation.participants_count < packetMaxLimit) {
        this.participantsCount = prepareInfo.data.conversation.participants_count
      } else {
        this.participantsCount = packetMaxLimit
      }
    }
    this.loading = false
  },
  computed: {
    validated () {
      if (this.form.amount && this.form.shares && this.selectedAsset && this.form.shares >= 1 && this.form.shares <= packetMaxLimit) {
        return true
      }
      return false
    }
  },
  methods: {
    async pay () {
      let payload = {
        amount: this.form.amount,
        total_count: parseInt(this.form.shares),
        greeting: this.form.memo,
        conversation_id: uuid.v4(),
        asset_id: this.selectedAsset.asset_id
      }

      this.loading = true
      let createResp = await this.GLOBAL.api.packet.create(payload)
      if (createResp.error) {
        this.loading = false
        Toast('Error')
        return
      }

      let pkt = createResp.data
      setTimeout(() => {
        this.waitForPayment(pkt.packet_id)
      }, 2000)
      window.location.href = `mixin://pay?recipient=${CLIENT_ID}&asset=${this.selectedAsset.asset_id}&amount=${this.form.amount}&trace=${pkt.packet_id}&memo=${encodeURIComponent(pkt.greeting)}`
    },
    onChangeAsset (ix) {
      this.selectedAsset = this.assets[ix]
      this.form.memo = this.$t('prepare_packet.default_memo', {symbol: this.selectedAsset.symbol})
    },
    async waitForPayment (packetId) {
      let resp = await this.GLOBAL.api.packet.show(packetId)
      if (resp.error) {
        setTimeout(() => { this.waitForPayment(packetId) }, 1500);
        return;
      }
      var data = resp.data;
      switch (data.state) {
        case 'INITIAL':
          setTimeout(() => { this.waitForPayment(packetId) }, 1500);
          break;
        case 'PAID':
        case 'EXPIRED':
        case 'REFUNDED':
          this.loading = false
          this.$router.push('/packets/' + packetId)
          break;
      }
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.prepare-packet-page {
  padding-top: 60px;
}
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
