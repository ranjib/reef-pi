import React from 'react'
import {connect} from 'react-redux'
import $ from 'jquery'
import { DropdownButton, MenuItem } from 'react-bootstrap'
import {fetchEquipments} from './redux/actions'

class selectEquipment extends React.Component {
  constructor (props) {
    super(props)
    var equipment = {id: props.active, name: ''}
    $.each(props.equipments, function (i, eq) {
      if (eq.id === equipment.id) {
        equipment = eq
      }
    })
    this.state = {
      equipment: equipment
    }
    this.equipmentList = this.equipmentList.bind(this)
    this.setEquipment = this.setEquipment.bind(this)
  }

  componentDidMount () {
    this.props.fetchEquipments()
  }

  equipmentList () {
    var menuItems = [ <MenuItem key='none' active={this.state.equipment === undefined} eventKey='none'>-</MenuItem> ]
    $.each(this.props.equipments, function (k, v) {
      var active = false
      if (this.state.equipment !== undefined) {
        active = this.state.equipment.id === v.id
      }
      menuItems.push(<MenuItem key={k} active={active} eventKey={k}><span id={this.props.id + '-' + v.name}>{v.name}</span></MenuItem>)
    }.bind(this))
    return menuItems
  }

  setEquipment (k, ev) {
    if (k === 'none') {
      this.setState({
        equipment: undefined
      })
      this.props.update('')
      return
    }
    var eq = this.props.equipments[k]
    this.setState({
      equipment: eq
    })
    this.props.update(eq.id)
  }

  render () {
    var readOnly = this.props.readOnly !== undefined ? this.props.readOnly : false
    var eqName = ''
    if (this.state.equipment !== undefined) {
      eqName = this.state.equipment.name
    }
    return (
      <div className='container'>
        <DropdownButton title={eqName} id={this.props.id} onSelect={this.setEquipment} disabled={readOnly}>
          {this.equipmentList()}
        </DropdownButton>
      </div>
    )
  }
}
const mapStateToProps = (state) => {
  return { equipments: state.equipments }
}

const mapDispatchToProps = (dispatch) => {
  return {fetchEquipments: () => dispatch(fetchEquipments())}
}

const SelectEquipment = connect(mapStateToProps, mapDispatchToProps)(selectEquipment)
export default SelectEquipment
