import fetchInlets from './actions'

export const rootReducer = (state, action) => {
  switch (action.type) {
    case 'INFO_LOADED':
      return { ...state, info: action.payload }
    case 'CAPABILITIES_LOADED':
      return { ...state, capabilities: action.payload }
    case 'JACKS_LOADED':
      return { ...state, jacks: action.payload }
    case 'INLETS_LOADED':
      return { ...state, inlets: action.payload }
    case 'EQUIPMENTS_LOADED':
      return { ...state, equipments: action.payload }
    case 'HEALTH_STATS_LOADED':
      return { ...state, health_stats: action.payload }
    case 'CREDS_UPDATED':
      return state
    default:
      return state
  }
}
