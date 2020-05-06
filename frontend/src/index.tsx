import {createElement, Context, Element} from "@bikeshaving/crank";
import {renderer} from "@bikeshaving/crank/dom";

import App from "./App";

function Index(this: Context): Element {
    return <App />;
};

renderer.render(<Index />, document.getElementById("root")!);