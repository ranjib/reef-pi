import React from 'react'
import Enzyme, {shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import Channel from './channel'
import Chart from './chart'
import Main from './main'
import LightForm from './light_form'
import Light from './light'
import AutoProfile from './auto_profile'
import DiurnalProfile from './diurnal_profile'
import FixedProfile from './fixed_profile'
import Profile from './profile' 
import Percent from './percent'
import configureMockStore from 'redux-mock-store'
import thunk from 'redux-thunk'
import 'isomorphic-fetch'
import {mockLocalStorage} from '../utils/test_helper'
window.localStorage = mockLocalStorage()

Enzyme.configure({ adapter: new Adapter() })
const mockStore = configureMockStore([thunk])

describe('Lighting ui', () => {
  const ev = {
    target: {value: 10}
  }
  const light = {
    id: '1',
    name: 'foo',
    jack: '1',
    channels: {
      '1': {
        pin: 0,
        color: '',
        profile: {
          type: 'auto',
          config: {values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]}
        }
      }
    }
  }
  it('<Main />', () => {
    const jacks = [{id: '1', name: 'foo'}]
    const m = shallow(<Main store={mockStore({lights: [light], jacks: jacks})} />).dive().instance()
    m.setJack(0, {})
    m.toggleAddLightDiv()
    m.addLight()
  })

  it('<LightForm />', () => {
    const m = shallow(<LightForm />).instance()
  })

  it('<Light />', () => {
    const values = { config: light }
    const m = shallow(<Light values={values} config={light} save={() => {}} remove={() => true} />).instance()
    m.toggleExpand()  
  })

  it('<Chart />', () => {
    shallow(<Chart store={mockStore({lights: [light]})} light_id='1' />).dive()
  })

  it('<Channel />', () => {
    const m = shallow(<Channel channel={light.channels['1']} onChangeHandler={() => {}} />)
    
  })

  it('<Profile /> fixed', () => {
    const wrapper = shallow(<Profile type='fixed' onChangeHandler={() => true} />)    
    expect(wrapper.find(FixedProfile).length).toBe(1)
    expect(wrapper.find(AutoProfile).length).toBe(0)
    expect(wrapper.find(DiurnalProfile).length).toBe(0)
  })
  
  it('<Profile /> auto', () => {
    const wrapper = shallow(<Profile type='auto' onChangeHandler={() => true} />)    
    expect(wrapper.find(FixedProfile).length).toBe(0)
    expect(wrapper.find(AutoProfile).length).toBe(1)
    expect(wrapper.find(DiurnalProfile).length).toBe(0)
  })
  
  it('<Profile /> diurnal', () => {
    const wrapper = shallow(<Profile type='diurnal' onChangeHandler={() => true} />)    
    expect(wrapper.find(FixedProfile).length).toBe(0)
    expect(wrapper.find(AutoProfile).length).toBe(0)
    expect(wrapper.find(DiurnalProfile).length).toBe(1)
  })

  it('<AutoProfile />', () => {
    const m = shallow(<AutoProfile onChangeHandler={() => true} />).instance()
    m.curry(1)(ev)
  })

  it('<DiurnalProfile />', () => {
    const m = shallow(<DiurnalProfile onChangeHandler={() => true} />).instance()
  })

  it('<FixedProfile />', () => {
    const m = shallow(<FixedProfile onChangeHandler={() => true} />).instance()
    m.handleChange(ev)
  })

  it('<Percent />', () => {
    const wrapper = shallow(<Percent value="4" onChange={() => true} />)
    wrapper.find('input').simulate('change', {target: {value: 34}})
  })

})
