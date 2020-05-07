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
    let autoscroll = true;

    const onFilterChange = (ev: Event) => {
        filter = (ev.currentTarget as HTMLInputElement).value;
        this.refresh();
    };

    const onAutoscrollChange = (ev: Event) => {
        autoscroll = (ev.currentTarget as HTMLInputElement).checked;
        this.refresh();
    }
    
    await new Promise((resolve) => webSocket.onopen = () => resolve());
    for await ({} of this) {
        let filteredData = data.filter(dataFilter(filter)).filter((d, i) => i < 20);
        yield (
            <div class={css(styles.page)}>
                <div class={css(styles.header)}>DevLog listening on tcp://localhost:9090/</div>
                <div class={css(styles.controls)}>
                    <input 
                        type="text"
                        name="search"
                        placeholder="Search..."
                        value={filter} 
                        oninput={onFilterChange} 
                        class={css(styles.searchInput)} 
                        />
                    <div class={css(styles.autoscroll)}>
                        <label
                            for="autoscroll"
                            class={css(styles.autoscrollLabel)}
                            >
                                Auto Scroll:
                        </label>
                        <input
                            type="checkbox"
                            name="autoscroll"
                            checked={autoscroll}
                            onchange={onAutoscrollChange}
                            class={css(styles.autoscrollCheckbox)}
                            />
                    </div>
                </div>
                <div class={css(styles.dataContainer)} id="dataContainer">
                    {filteredData.map((d, i) => (
                        <JsonView crank-key={d.id} data={d} />
                    ))}
                </div>
            </div>
        );
        if (autoscroll) {
            let container = document.getElementById("dataContainer");
            if (container)
                container.scrollTop = container.scrollHeight;
        }
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
    controls: {
        display: 'flex',
    },
    autoscroll: {
        fontFamily: 'Monospace',
        marginLeft: '0.5em',
        whiteSpace: 'nowrap',
    },
    autoscrollLabel: {
        marginRight: '0.5em',
    },
    autoscrollCheckbox: {
        verticalAlign: 'bottom',
        position: 'relative',
        top: '-1px',
        overflow: 'hidden',
        margin: 0,
        padding: 0,
        width: '13px',
        height: '13px',
    },
    searchInput: {
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