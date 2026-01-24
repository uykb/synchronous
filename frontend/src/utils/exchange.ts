export interface ExchangeInfo {
  id: string
  name: string
  icon: string
  fields: ExchangeField[]
}

export interface ExchangeField {
  key: string
  label: string
  type: 'password' | 'text'
  placeholder: string
}

export const EXCHANGES: ExchangeInfo[] = [
  {
    id: 'binance',
    name: 'Binance',
    icon: 'https://cryptologos.cc/logos/binance-coin-bnb-logo.svg',
    fields: [
      { key: 'api_key', label: 'API Key', type: 'password', placeholder: '请输入 API Key' },
      { key: 'api_secret', label: 'Secret Key', type: 'password', placeholder: '请输入 Secret Key' }
    ]
  },
  {
    id: 'okx',
    name: 'OKX',
    icon: 'https://cryptologos.cc/logos/okb-okb-logo.svg',
    fields: [
      { key: 'api_key', label: 'API Key', type: 'password', placeholder: '请输入 API Key' },
      { key: 'api_secret', label: 'Secret Key', type: 'password', placeholder: '请输入 Secret Key' },
      { key: 'passphrase', label: 'Passphrase', type: 'password', placeholder: '请输入 Passphrase' }
    ]
  },
  {
    id: 'bybit',
    name: 'Bybit',
    icon: 'https://s2.coinmarketcap.com/static/img/exchanges/64x64/521.png',
    fields: [
      { key: 'api_key', label: 'API Key', type: 'password', placeholder: '请输入 API Key' },
      { key: 'api_secret', label: 'Secret Key', type: 'password', placeholder: '请输入 Secret Key' }
    ]
  },
  {
    id: 'backpack',
    name: 'Backpack',
    icon: 'https://avatars.githubusercontent.com/u/112170256?s=200&v=4',
    fields: [
      { key: 'api_key', label: 'API Key', type: 'password', placeholder: '请输入 API Key (公钥)' },
      { key: 'api_secret', label: 'Secret Key', type: 'password', placeholder: '请输入 Secret Key (Ed25519 私钥)' }
    ]
  },
  {
    id: 'lighter',
    name: 'Lighter',
    icon: 'https://framerusercontent.com/images/PwXZowDeTtB9awXKHOIK8lnD4.png',
    fields: [
      { key: 'api_key', label: 'API Key', type: 'password', placeholder: '请输入 API 公钥' },
      { key: 'api_secret', label: 'API Secret', type: 'password', placeholder: '请输入 API 私钥' }
    ]
  }
]

export function getExchangeById(id: string): ExchangeInfo | undefined {
  return EXCHANGES.find(e => e.id === id)
}
