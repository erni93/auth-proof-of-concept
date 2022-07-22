import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { faCompass } from '@fortawesome/free-solid-svg-icons';

import { routes } from 'src/app/app-routing.module';

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss'],
})
export class NavigationComponent {
  public faCompass = faCompass;
  public activeRoutes: string[] = [];

  constructor(router: Router) {
    this.activeRoutes = this.getNavigationRoutes(router.url.replace('/', ''));
  }

  private getNavigationRoutes(actualRoute: string): string[] {
    console.log('route', actualRoute);
    return routes.reduce((acc, r) => {
      if (r.path && r.path !== actualRoute) {
        acc.push(r.path);
      }
      return acc;
    }, [] as string[]);
  }
}
