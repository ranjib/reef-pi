import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import ATO from './ato'
import New from './new'
import Main from './main'
import Chart from './chart'
import configureMockStore from 'redux-mock-store'
import thunk from 'redux-thunk'
import 'isomorphic-fetch'

Enzyme.configure({ adapter: new Adapter() })
const mockStore = configureMockStore([thunk])

describe('ATO ui', () => {
  const state = {
    ato_usage: { '1': {} },
    atos: [{ name: 'foo', id: '1' }]
  }
  it('<ATO />', () => {
    const n = shallow(<ATO store={mockStore()} data={{ id: '1', period: 10 }} />)
    const m = n.dive().instance()
    m.save()
    m.state.readOnly = false
    m.setInlet('1')
    m.update('period')({ target: { value: 'abc' } })
    m.update('period')({ target: { value: 10 } })
    m.updateCheckBox('control')({ target: {} })
    m.updatePump('1')
    m.save()
    m.remove()
    m.expand()
    m.detailsUI()
    const z = shallow(<ATO store={mockStore({ readOnly: false })} data={{ id: '1', period: 'foo', control: true }} />)
      .dive()
      .instance()
    z.state.readOnly = false
    z.state.ato = { id: '1', period: 'foo', control: true }
    z.save()
    z.showControl()
  })

  it('<New />', () => {
    const m = shallow(<New store={mockStore()} />)
      .dive()
      .instance()
    m.toggle()
    m.update('name')({ target: { value: 's' } })
    m.updateEnable({ target: { checked: true } })
    m.setInlet('1')
    m.add()
  })

  it('<Main />', () => {
    shallow(<Main store={mockStore({ atos: [{ id: '1', period: 10 }] })} />).dive()
  })

  it('<Chart />', () => {
    shallow(<Chart ato_id='1 ' store={mockStore(state)} />)
      .dive()
      .instance()
  })
})
