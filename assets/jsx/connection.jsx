import React from 'react';
import $ from 'jquery';
import { DropdownButton, MenuItem } from 'react-bootstrap';

export default class Connection extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
          boards: [],
          board: 0,
          pin: 0
        }
        this.fetchBoards = this.fetchBoards.bind(this);
        this.boardList = this.boardList.bind(this);
        this.pinList = this.pinList.bind(this);
        this.setBoard = this.setBoard.bind(this);
        this.setPin = this.setPin.bind(this);
    }

    componentWillMount(){
      this.fetchBoards();
    }

    fetchBoards(){
        $.ajax({
            url: '/api/boards',
            type: 'GET',
            dataType: 'json',
            success: function(data) {
              this.setState({
                boards: data
              });
            }.bind(this),
            error: function(xhr, status, err) {
              console.log(err.toString());
            }.bind(this)
        });
    }

    setBoard(k, ev){
      this.setState({
        board: Number(k),
      });
      var payload ={
        board: this.state.boards[k].id,
        pin: this.state.pin
      }
      this.props.updateHook(payload);
    }

    setPin(k, ev){
      this.setState({
        pin: Number(k)
      });

      var payload ={
        board: this.state.boards[this.state.board].id,
        pin: this.state.pin
      }
      this.props.updateHook(payload);
    }

    boardList(){
      var list = []
      $.each(this.state.boards, function(k, v){
        list.push(<MenuItem key={v.name} eventKey={k}> {v.name}</MenuItem>);
      }.bind(this));
      return list;
    }

    pinList(){
      var list = []
      if(this.state.boards.length == 0){
        return list;
      }
      var l = this.state.boards[this.state.board].pins
      for(var i =1; i<= l; i++) {
        list.push(<MenuItem key={i} eventKey={i}> {i}</MenuItem>);
      }
      return list;
    }


    render() {
     var boardName = ""
     if(this.state.boards[this.state.board] != undefined){
       boardName = this.state.boards[this.state.board].name
     }
      return(
          <div>
            <DropdownButton  title={boardName} id="board-name" onSelect={this.setBoard}>
              {this.boardList()}
            </DropdownButton>
            <DropdownButton  title={this.state.pin} id="pin-number" onSelect={this.setPin}>
              {this.pinList()}
            </DropdownButton>
          </div>
          );
    }
}
