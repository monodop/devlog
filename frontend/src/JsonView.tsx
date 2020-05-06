import { createElement, Fragment, Props } from "@bikeshaving/crank";
import { Context } from "@bikeshaving/crank";
import { css, StyleSheet } from 'aphrodite';

export interface IJsonViewProps<T> {
    data: T,
};

type SMap<TValue> = {
    [key: string]: TValue;
}

function AnyView(this: Context<IJsonViewProps<any>>, {data}: IJsonViewProps<any>) {
    if (typeof data === 'object') {
        return <ObjectView data={data} />;
    }
    else if (typeof data === 'string' || typeof data === 'number') {
        return <StringView data={data} />;
    }
    else {
        return <div>unknown</div>
    }
}

function ObjectView(this: Context<IJsonViewProps<SMap<any>>>, {data}: IJsonViewProps<SMap<any>>) {
    return (
        <div class={css(styles.objectview)}>
            {Object.keys(data).map(k => (
                <Fragment crank-key={k}>
                    <div class={css(styles.objectviewKey)}>
                        {k}
                    </div>
                    <div class={css(styles.objectviewValue)}>
                        <AnyView data={data[k]} />
                    </div>
                </Fragment>
            ))}
        </div>
    )
}

function StringView(this: Context<IJsonViewProps<string|number>>, {data}: IJsonViewProps<string|number>) {
    return (
        <span>{data}</span>
    )
}

export default function JsonView(this: Context<IJsonViewProps<any>>, {data}: IJsonViewProps<any>) {
    
    return (
        <section class={css(styles.jsonview)}>
            <AnyView data={data} />
        </section>
    )
}

const styles = StyleSheet.create({
    jsonview: {
        borderBottom: '1px solid #eee',
        // margin: '0.5em 0',
        padding: '0.5em',
        // backgroundColor: '#efefef',
        fontFamily: 'Monospace',
    },
    objectview: {
        display: 'grid',
        gridTemplateColumns: '0fr 1fr',
        gridRowGap: '0.5em',
        gridColumnGap: '1em',
    },
    objectviewKey: {
        // padding: '0.5em 0.25em',
        // paddingRight: '0.5em',
        color: '#555',
    },
    objectviewValue: {
        // padding: '0.5em 0.25em',
    },
});