import { h, Component } from "preact";
import { TimeSelection } from "../models/TimeSelection";

enum Selection {
    LAST_15_MINUTES,
    LAST_60_MINUTES,
    LAST_4_HOURS,
    LAST_24_HOURS,
    LAST_7_DAYS,
    LAST_30_DAYS
}

interface Option {
    value: Selection;
    name: string;
    ts: TimeSelection;
}

const options: Option[] = [
    { value: Selection.LAST_15_MINUTES, name: 'Last 15 minutes', ts: { relativeTime: '-15m' } },
    { value: Selection.LAST_60_MINUTES, name: 'Last 60 minutes', ts: { relativeTime: '-60m' } },
    { value: Selection.LAST_4_HOURS, name: 'Last 4 hours', ts: { relativeTime: '-4h' } },
    { value: Selection.LAST_24_HOURS, name: 'Last 24 hours', ts: { relativeTime: '-24h' } },
    { value: Selection.LAST_7_DAYS, name: 'Last 7 days', ts: { relativeTime: '-168h' } },
    { value: Selection.LAST_30_DAYS, name: 'Last 30 days', ts: { relativeTime: '-720h' } },
];

interface TimeSelectProps {
    onTimeSelected: (newTime: TimeSelection) => void;
}

interface TimeSelectState {
    currentSelection: Selection;
    selectedOptionName: string;
}

export class TimeSelect extends Component<TimeSelectProps, TimeSelectState> {
    constructor(props: TimeSelectProps) {
        super(props);
        this.state = {
            currentSelection: options[0].value,
            selectedOptionName: options[0].name,
        };
        this.props.onTimeSelected(options[0].ts)
    }

    render() {
        return [
            <button
                class="btn btn-outline-secondary dropdown-toggle"
                type="button"
                data-toggle="dropdown"
                aria-haspopup="true"
                aria-expanded="false">
                {this.state.selectedOptionName}
            </button>,
            <div class="dropdown-menu dropdown-menu-right" style="min-width: 276px;">
                {options.map(o =>
                    <button
                        type="button"
                        class={"dropdown-item" + (this.state.currentSelection === o.value ? " bg-info text-white" : "")}
                        onClick={() => this.onSelection(o)}>
                        {o.name}
                    </button>
                )}
            </div>
        ];
    }

    private onSelection(o: Option) {
        this.setState({
            currentSelection: o.value,
            selectedOptionName: o.name
        });
        this.props.onTimeSelected(o.ts);
    }
}