import React, { Component } from 'react';
import '../assets/css/main.css';
import RouterWrap from '../router';

export default class App extends Component {
  render() {
    return (
      <div className="App">
        <RouterWrap />
      </div>
    );
  }
}