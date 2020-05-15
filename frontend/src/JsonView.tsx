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
    if (typeof data === 'string') {
        return <StringView data={data} />;
    }
    else if (typeof data === 'number') {
        return <NumberView data={data} />;
    }
    else if (typeof data === 'boolean') {
        return <BooleanView data={data} />;
    }
    else if (data === null) {
        return <NullView data={data} />;
    }
    else if (typeof data === 'object') {
        return <ObjectView data={data} />;
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

function StringView(this: Context<IJsonViewProps<string>>, {data}: IJsonViewProps<string>) {
    return (
        <span class={css(styles.string)}>{data}</span>
    )
}

function NumberView(this: Context<IJsonViewProps<number>>, {data}: IJsonViewProps<number>) {
    return (
        <span class={css(styles.number)}>{data}</span>
    )
}

function BooleanView(this: Context<null>, {data}: IJsonViewProps<boolean>) {
    return (
        <span class={css(styles.boolean)}>{data ? 'true' : 'false'}</span>
    )
}

function NullView(this: Context<null>, {data}: IJsonViewProps<null>) {
    return (
        <span class={css(styles.null)}>null</span>
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

    string: {
        color: '#090',
    },
    number: {
        color: '#33d',
    },
    boolean: {
        color: '#33d',
    },
    null: {
        color: '#d00',
    },
});