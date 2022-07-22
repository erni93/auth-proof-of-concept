import { Component } from '@angular/core';
import { faUserClock } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['../shared/screens-shared.scss', './home.component.scss'],
})
export class HomeComponent {
  faUserClock = faUserClock;
}
