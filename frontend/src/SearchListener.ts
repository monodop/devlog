import { ANTLRInputStream, CommonTokenStream } from "antlr4ts";
import { FilterGrammarLexer } from "./FilterGrammar/FilterGrammarLexer";
import { FilterGrammarParser, PlaintextContext, Filter_simpleContext, Filter_keyContext } from "./FilterGrammar/FilterGrammarParser";
import { ParseTreeListener } from "antlr4ts/tree/ParseTreeListener";
import { ParseTreeWalker } from "antlr4ts/tree/ParseTreeWalker";
import { FilterGrammarListener } from "./FilterGrammar/FilterGrammarListener";


export function dataFilter(f: string): (data: any) => boolean {
    return (data) => {

        const inputStream = new ANTLRInputStream(f);
        const lexer = new FilterGrammarLexer(inputStream);
        const tokenStream = new CommonTokenStream(lexer);
        const parser = new FilterGrammarParser(tokenStream);

        let tree = parser.entry();
        let listener = new SearchListener(data);
        ParseTreeWalker.DEFAULT.walk(listener as ParseTreeListener, tree);
        return listener.matches;
    }
}

class SearchListener implements FilterGrammarListener {

    private _data: any;
    matches: boolean = true;

    constructor(data: any) {
        this._data = data;
    }

    getPlaintextText(context: PlaintextContext): string {
        if (context.WORD()) {
            return context.WORD()!.text;
        }
        else if(context.STRING()) {
            let word = context.STRING()!.text;
            return word.substr(1, word.length - 2)
                .replace('\\"', '"')
                .replace('\\\\', '\\');
        }
        return '';
    }

    enterFilter_simple(context: Filter_simpleContext) {
        let plaintext = context.plaintext();
        if (plaintext) {
            let word: string = this.getPlaintextText(plaintext);
            if (word && !plaintextSearch(this._data, word)) {
                this.matches = false;
            }
        }
    }

    enterFilter_key(context: Filter_keyContext) {
        let equals = context.equals();
        if (equals) {
            let key = this.getPlaintextText(equals.plaintext()[0]);
            let value = this.getPlaintextText(equals.plaintext()[1]).toString().toLowerCase();
            let compareTo = (this._data[key] || '').toString().toLowerCase();

            let startsWith = false;
            let endsWith = false;

            if (value.startsWith("*")) {
                endsWith = true;
                value = value.substr(1);
            }
            if (value.endsWith("*")) {
                startsWith = true;
                value = value.substr(0, value.length - 1);
            }

            if (startsWith && endsWith) {
                if (compareTo.indexOf(value) < 0) {
                    this.matches = false;
                }
            }
            else if (startsWith) {
                if (!compareTo.startsWith(value)) {
                    this.matches = false;
                }
            }
            else if (endsWith) {
                if (!compareTo.endsWith(value)) {
                    this.matches = false;
                }
            }
            else {
                if (compareTo !== value) {
                    this.matches = false;
                }
            }

        }
    }
}

function plaintextSearch(data: any, f: string): boolean {
    if (data === null || data === undefined) {
        return false;
    }
    else if (typeof data === 'object') {
        return Object.keys(data).some(k => plaintextSearch(data[k], f));
    }
    else if (typeof data === 'string' || typeof data === 'number') {
        return data.toString().toLowerCase().indexOf(f.toString().toLowerCase()) >= 0;
    }
    return false;
}