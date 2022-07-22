import { Component } from '@angular/core';

import { faStopwatch } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-sessions',
  templateUrl: './sessions.component.html',
  styleUrls: ['../shared/screens-shared.scss'],
})
export class SessionsComponent {
  public faStopwatch = faStopwatch;
}
