import { createElement, Fragment } from "@bikeshaving/crank";
import { Context } from "@bikeshaving/crank";
import JsonView from "./JsonView";
import { dataFilter } from "./SearchListener";
import { StyleSheet, css } from 'aphrodite';

export interface IAppProps {};

export default async function* App(this: Context<IAppProps>, {}: IAppProps) {
    const data: any[] = [];
    const webSocket = new WebSocket(`ws://${location.hostname}:9091/ws`);
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
            <div class={css(styles.page)}>
                <div class={css(styles.header)}>DevLog listening on tcp://localhost:9090/</div>
                <input 
                    type="text"
                    placeholder="Search..."
                    value={filter} 
                    oninput={onFilterChange} 
                    class={css(styles.input)} 
                    />
                <div class={css(styles.dataContainer)}>
                    {filteredData.map((d, i) => (
                        <JsonView crank-key={d.id} data={d} />
                    ))}
                </div>
            </div>
        );
    }
}

const styles = StyleSheet.create({
    page: {
        height: '100vh',
        boxSizing: 'border-box',
        padding: '0.5em',
        display: 'flex',
        flexDirection: 'column',
    },
    input: {
        boxSizing: 'border-box',
        width: '100%',
        marginBottom: '0.5em',
    },
    header: {
        fontFamily: 'Monospace',
        backgroundColor: '#ccc',
        padding: '1em',
        marginBottom: '1em',
    },
    dataContainer: {
        flex: 1,
        overflowY: 'scroll',
    },
});