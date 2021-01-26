import React,{Component} from 'react'; 
import {BrowserRouter,Route,Switch} from 'react-router-dom';
import AdminContent from './layouts/AdminLayout';
import LoginContent from './layouts/LoginLayout';
import {Exception404, Exception500} from "./layouts/ExceptionContent";

export default class RouterWrap extends Component{ 
    
    render(){ 
        return ( 
            <div id="router" style={{backgroundColor:"#f0f2f5"}}>
                <BrowserRouter>
                    <Switch>
                        <Route path="/" component={LoginContent} exact />
                        <Route path="/admin" component={AdminContent} />
                        <Route path="/login" component={LoginContent} />
                        <Route component={Exception404} />
                        <Route path="/exception500" component={Exception500} />
                    </Switch>
                </BrowserRouter>
            </div> 
        )
    } 

}
