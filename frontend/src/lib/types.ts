export type LoadStatus = 'Tendered' | 'Covered' | string

export interface PartyDetails {
  externalTMSId: string
  name: string
  addressLine1: string
  addressLine2: string
  city: string
  state: string
  zipcode: string
  country: string
  contact: string
  phone: string
  email: string
  refNumber?: string
}

export interface PickupDetails extends PartyDetails {
  businessHours: string
  refNumber: string
  readyTime: string
  apptTime: string
  apptNote: string
  timezone: string
  warehouseId: string
}

export interface ConsigneeDetails extends PartyDetails {
  businessHours: string
  refNumber: string
  mustDeliver: string
  apptTime: string
  apptNote: string
  timezone: string
  warehouseId: string
}

export interface CarrierDetails {
  mcNumber: string
  dotNumber: string
  name: string
  phone: string
  dispatcher: string
  sealNumber: string
  scac: string
  firstDriverName: string
  firstDriverPhone: string
  secondDriverName: string
  secondDriverPhone: string
  email: string
  dispatchCity: string
  dispatchState: string
  externalTMSTruckId: string
  externalTMSTrailerId: string
  confirmationSentTime: string
  confirmationReceivedTime: string
  dispatchedTime: string
  expectedPickupTime: string
  pickupStart: string
  pickupEnd: string
  expectedDeliveryTime: string
  deliveryStart: string
  deliveryEnd: string
  signedBy: string
  externalTMSId: string
}

export interface RateData {
  customerRateType: string
  customerNumHours: number
  customerLhRateUsd: number
  fscPercent: number
  fscPerMile: number
  carrierRateType: string
  carrierNumHours: number
  carrierLhRateUsd: number
  carrierMaxRate: number
  netProfitUsd: number
  profitPercent: number
}

export interface Specifications {
  minTempFahrenheit: number
  maxTempFahrenheit: number
  liftgatePickup: boolean
  liftgateDelivery: boolean
  insidePickup: boolean
  insideDelivery: boolean
  tarps: boolean
  oversized: boolean
  hazmat: boolean
  straps: boolean
  permits: boolean
  escorts: boolean
  seal: boolean
  customBonded: boolean
  labor: boolean
}

export interface Load {
  externalTMSLoadID: string
  freightLoadID: string
  status: LoadStatus
  customer: PartyDetails
  billTo: Omit<PartyDetails, 'refNumber'>
  pickup: PickupDetails
  consignee: ConsigneeDetails
  carrier: CarrierDetails
  rateData: RateData
  specifications: Specifications
  inPalletCount: number
  outPalletCount: number
  numCommodities: number
  totalWeight: number
  billableWeight: number
  poNums: string
  operator: string
  routeMiles: number
}

export interface LoadsResponse {
  data: Load[]
  pagination: {
    total: number
    pages: number
    page: number
    limit: number
    approximate: boolean
  }
}

export interface CreateLoadResponse {
  id: number
  createdAt: string
}

export interface ListLoadsParams {
  page: number
  limit: number
  status?: string
  customerId?: string
  pickupDateSearchFrom?: string
  pickupDateSearchTo?: string
}
