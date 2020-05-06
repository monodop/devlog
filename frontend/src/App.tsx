import { createElement, Fragment } from "@bikeshaving/crank";
import { Context } from "@bikeshaving/crank";
import JsonView from "./JsonView";
import { dataFilter } from "./SearchListener";
import { StyleSheet, css } from 'aphrodite';

export interface IAppProps {};

export default async function* App(this: Context<IAppProps>, {}: IAppProps) {
    const data: any[] = [];
    const webSocket = new WebSocket("ws://localhost:9091/ws");
    webSocket.onmessage = (ev) => {
        const parsed = JSON.parse(JSON.parse(ev.data).Message);
        console.log(parsed);
        data.push(parsed);
        while (data.length > 10000) {
            data.shift();
        }
        this.refresh();
    }

    let filter = '';

    const onFilterChange = (ev: Event) => {
        filter = (ev.currentTarget as HTMLInputElement).value;
        this.refresh();
    };
    
    await new Promise((resolve) => webSocket.onopen = () => resolve());
    for await ({} of this) {
        let filteredData = data.filter(dataFilter(filter)).filter((d, i) => i < 20);
        yield (
            <Fragment>
                <div class={css(styles.header)}>DevLog listening on tcp://localhost:9090/</div>
                <input 
                    type="text"
                    placeholder="Search..."
                    value={filter} 
                    oninput={onFilterChange} 
                    class={css(styles.input)} 
                    />
                <div>
                    {filteredData.map((d, i) => (
                        <JsonView crank-key={d.id} data={d} />
                    ))}
                </div>
            </Fragment>
        );
    }
}

const styles = StyleSheet.create({
    input: {
        boxSizing: 'border-box',
        width: '100%',
    },
    header: {
        fontFamily: 'Monospace',
        backgroundColor: '#ccc',
        padding: '1em',
        marginBottom: '1em',
    },
});