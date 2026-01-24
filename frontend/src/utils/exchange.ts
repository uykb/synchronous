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
    icon: 'https://cryptologos.cc/logos/bitcoin-btc-logo.svg',
    fields: [
      { key: 'api_key', label: 'API Key', type: 'password', placeholder: '请输入 API Key' },
      { key: 'api_secret', label: 'Secret Key', type: 'password', placeholder: '请输入 Secret Key' }
    ]
  }
]

export function getExchangeById(id: string): ExchangeInfo | undefined {
  return EXCHANGES.find(e => e.id === id)
}
