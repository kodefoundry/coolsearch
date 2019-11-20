import React, {Component} from "react";
import Resultpage from './Resultpage';
import {subscribeToSearchResult} from '../Socketclient';


function uuidv4() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        var r = Math.random() * 16 | 0,
        v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

class Homepage extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            items: [],
            currentItem: {
                text: '',
                key: '',
                status:''
            },
            uuid: ''
        }
        this.addItem = this.addItem.bind(this);
        this.handleInput = this.handleInput.bind(this);
        this.deleteItem = this.deleteItem.bind(this);
        this.setUpdate = this.setUpdate.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
        
        subscribeToSearchResult((err, uuid) => {
          this.setState({uuid: uuid})
        });
        
    }

    handleSubmit(uuid, searchword) {
        const url = "http://localhost:8082/v1/api/postquery";
        const encodedSearchword = encodeURI(searchword);
        const data = {
            uuid: uuid,
            searchword: encodedSearchword
        };
        console.log("Body : " + JSON.stringify(data));
        fetch(url, {
            mode: 'cors',
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        }).then(res => res.json()).catch(error => console.error('Error:', error)).then(response => console.log('Success:', response));
    }    

    addItem(e) {
        e.preventDefault();
        const newItem = this.state.currentItem;
        if (newItem.text !== "") {
            newItem.key = uuidv4()
            newItem.status = 'Pending'
            console.log("New search:",this.state.currentItem.text,"UUID: ", this.state.currentItem.key);
            const items = [
                newItem,
                ...this.state.items
            ];
            this.setState({
                items: items,
                currentItem: {
                    text: '',
                    key: '',
                    status:''
                }
            })
          //handleSubmit(newItem.key, newItem.text)
            const url = "http://localhost:8082/api/postquery";
            const encodedSearchword = encodeURI(newItem.text);
            const data = {
                uuid: newItem.key,
                searchword: encodedSearchword
            };
            console.log("Body : " + JSON.stringify(data));
            fetch(url, {
                mode: 'cors',
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            }).then(res => res.json()).catch(error => console.error('Error:', error)).then(response => console.log('Success:', response));
        }
    }

    handleInput(e) {
        this.setState({
            currentItem: {
                text: e.target.value,
                key: '',
                status: ''
            }
        })
    }

    deleteItem(key) {
        const filteredItems = this.state.items.filter(item => item.key !== key);
        this.setState({items: filteredItems})

    }

    setUpdate(text, key) {
        console.log("items:" + this.state.items);
        const items = this.state.items;
        items.map(item => {
            if (item.key === key) {
                console.log(item.key + "    " + key)
                item.text = text;
            }
        })
        this.setState({items: items})
    }

    render() {
        return (
            <div>
                <header>
                    <form id="to-do-form" onSubmit={this.addItem} className="App" >
                        <input type="text" placeholder="Enter search" value={ this.state.currentItem.text } onChange={this.handleInput}></input>
                        <button type="submit">Search</button>
                    </form>
                    <Resultpage availableuuid={this.state.uuid} items={this.state.items} deleteItem={this.deleteItem} setUpdate={this.setUpdate}/>
                </header>
            </div>
        )
    }
}

export default Homepage;
