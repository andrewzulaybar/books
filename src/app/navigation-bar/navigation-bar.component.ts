import { Component, AfterViewInit, ElementRef, ViewChild } from '@angular/core';

@Component({
  selector: 'app-navigation-bar',
  templateUrl: './navigation-bar.component.html',
  styleUrls: ['./navigation-bar.component.css'],
})
export class NavigationBarComponent implements AfterViewInit {
  @ViewChild('navBar', { read: ElementRef })
  private navigationBar: ElementRef;

  public get element(): HTMLElement {
    return this.navigationBar.nativeElement;
  }

  public ngAfterViewInit() {
    this.addBoxShadow();
    window.addEventListener('scroll', () => this.addBoxShadow());
  }

  private addBoxShadow(): void {
    const opacity: string = this.calculateBoxShadowOpacity().toString();
    this.element.style.setProperty('--box-shadow-opacity', opacity);
  }

  private calculateBoxShadowOpacity(): number {
    const { scrollTop }: { scrollTop: number } = document.documentElement;
    const quarterViewport: number = window.innerHeight / 4;
    return Math.min(scrollTop / quarterViewport, 1);
  }
}
