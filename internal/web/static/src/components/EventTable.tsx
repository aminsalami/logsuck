/**
 * Copyright 2021 The Logsuck Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { h } from "preact";
import { LogEvent } from "../models/Event";

export interface EventTableProps {
  events: LogEvent[];
}

export const EventTable = ({ events }: EventTableProps) => (
  <table class="table table-hover search-result-table">
    <thead>
      <tr>
        <th scope="col">Time</th>
        <th scope="col">Event</th>
      </tr>
    </thead>
    <tbody>
      {events.map((e) => (
        <tr key={e.raw}>
          <td class="event-timestamp">{e.timestamp.toLocaleString()}</td>
          <td>
            <div
              style={{
                display: "flex",
                flexDirection: "column",
              }}
            >
              <div class="event-raw">{e.raw}</div>
              <hr
                style={{
                  width: "100%",
                  marginTop: "0.75rem",
                  marginBottom: "0.5rem",
                }}
              />
              <div class="event-additional">
                <dl class="row no-gutters" style={{ marginBottom: 0 }}>
                  <dt class="col-1">source</dt>
                  <dd class="col-1">{e.source}</dd>
                </dl>
              </div>
            </div>
          </td>
        </tr>
      ))}
    </tbody>
  </table>
);
